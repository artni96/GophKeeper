package user

import (
	"time"

	"github.com/google/uuid"
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserCreate struct {
	ID             *uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username       string     `gorm:"column:username"`
	HashedPassword string     `gorm:"column:hashed_password"`
}

type LoginRequest struct {
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}

type LoginResponse struct {
	Token string
	Keys  []UserKey
}

type User struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;"`
	Username       string    `gorm:"column:username"`
	HashedPassword string    `gorm:"column:hashed_password"`
}

type UserKeyCreate struct {
	UserID       uuid.UUID `gorm:"column:user_id;type:uuid;"`
	EncryptedKey []byte    `gorm:"column:encrypted_key"`
	Salt         []byte    `gorm:"column:salt"`
	IsActive     bool      `gorm:"column:is_active"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}

type UserKey struct {
	ID           uuid.UUID `gorm:"column:id;type:uuid;"`
	UserID       uuid.UUID `gorm:"column:user_id;type:uuid;"`
	EncryptedKey []byte    `gorm:"column:encrypted_key"`
	KeyID        uint64    `gorm:"column:key_id"`
	Salt         []byte    `gorm:"column:salt"`
	IsActive     bool      `gorm:"column:is_active"`
	CreatedAt    time.Time `gorm:"column:created_at"`
}
