package user

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/artni96/GophKeeper/internal/common/utils"
	"github.com/artni96/GophKeeper/internal/server/config"
	usermodel "github.com/artni96/GophKeeper/internal/server/model/user"
	"github.com/artni96/GophKeeper/internal/server/repository/user"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

var (
	ErrWrongUserOrPassword = errors.New("wrong user or password")
	ErrFailedToCreateUser  = errors.New("failed to create user")
)

// ServiceI implements User service interface.
type ServiceI interface {
	Create(ctx context.Context, entity usermodel.UserCreateRequest) error
	Login(ctx context.Context, entity usermodel.LoginRequest) (usermodel.LoginResponse, error)
}

// Service implements the User service business logic.
type Service struct {
	repo   user.RepositoryI
	cfg    *config.Config
	logger *zap.Logger
}

// NewService initializes and returns the new User service instance.
func NewService(cfg *config.Config, logger *zap.Logger, repo user.RepositoryI) *Service {
	return &Service{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

// Create creates a new User entity.
func (s *Service) Create(ctx context.Context, entity usermodel.UserCreateRequest) error {
	hashedPassword, err := s.hashPassword(entity.Password)
	if err != nil {
		s.logger.Info("failed to hash password", zap.Error(err))
		return err
	}

	userToCreate := usermodel.UserCreate{
		Username:       entity.Username,
		HashedPassword: hashedPassword,
	}

	salt, err := utils.GenerateRandomBytes(aes.BlockSize)
	if err != nil {
		s.logger.Info(
			"failed to generate salt for new user", zap.String("login", entity.Username), zap.Error(err))
		return ErrFailedToCreateUser
	}

	derivedKey := pbkdf2.Key([]byte(entity.Password), salt, 10000, 32, sha256.New)

	aeskey, err := utils.GenerateRandomBytes(2 * aes.BlockSize)
	if err != nil {
		s.logger.Info(
			"failed to generate key for new user", zap.String("login", entity.Username), zap.Error(err))
		return ErrFailedToCreateUser
	}

	aesblock, err := aes.NewCipher(derivedKey)
	if err != nil {
		s.logger.Info(
			"failed to new AES cipher for new user", zap.String("login", entity.Username), zap.Error(err))
		return ErrFailedToCreateUser
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		s.logger.Info("failed to new GCM", zap.String("login", entity.Username), zap.Error(err))
		return ErrFailedToCreateUser
	}

	nonce, err := utils.GenerateRandomBytes(aesgcm.NonceSize())
	if err != nil {
		s.logger.Info(
			"failed to generate nonce for new user", zap.String("login", entity.Username), zap.Error(err))
		return ErrFailedToCreateUser
	}

	ciphertext := aesgcm.Seal(nil, nonce, aeskey, nil)
	encryptedKey := append(nonce, ciphertext...)

	userKeyToCreate := usermodel.UserKeyCreate{
		EncryptedKey: encryptedKey,
		Salt:         salt,
		IsActive:     true,
		CreatedAt:    time.Now(),
	}

	err = s.repo.Create(ctx, userToCreate, userKeyToCreate)
	if err != nil {
		return err
	}
	return nil
}

// Login returns a jwt token for a user.
func (s *Service) Login(ctx context.Context, entity usermodel.LoginRequest) (usermodel.LoginResponse, error) {
	resp := usermodel.LoginResponse{}
	dbEntity, err := s.repo.GetByUsername(ctx, entity.Username)
	if err != nil {
		return resp, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbEntity.HashedPassword), []byte(entity.Password))
	if err != nil {
		s.logger.Info("wrong password", zap.Error(err))
		return resp, ErrWrongUserOrPassword
	}
	token, err := s.BuildJWTString(dbEntity.ID, s.cfg)
	if err != nil {
		s.logger.Info("failed to build JWT token", zap.Error(err))
		return resp, err
	}
	userKeys, err := s.repo.GetUserKeysList(ctx, dbEntity.ID)
	if err != nil {
		return resp, err
	}
	resp.Token = token
	resp.Keys = userKeys
	return resp, nil
}

// hashPassword hashes a plain password by bcrypt.
func (s *Service) hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hash), nil
}

// Claims represents claims for the JWT token.
type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

// BuildJWTString creates a JWT token by user's ID.
func (s *Service) BuildJWTString(userID uuid.UUID, cfg *config.Config) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.TokenExp)),
		},
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(cfg.SecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}
	return tokenString, nil
}

// GetUserIDFromJWT extracts the user ID from the JWT token.
func GetUserIDFromJWT(tokenString string, cfg *config.Config) uuid.UUID {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.SecretKey), nil
	})
	if err != nil {
		return uuid.Nil
	}

	if !token.Valid {
		return uuid.Nil
	}
	return claims.UserID
}
