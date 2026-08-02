package login

import (
	"context"
	"errors"
	"time"

	pb "github.com/artni96/GophKeeper/api/proto/logins"
	"github.com/artni96/GophKeeper/internal/interceptors"
	loginmodel "github.com/artni96/GophKeeper/internal/model/login"
	userrepo "github.com/artni96/GophKeeper/internal/repository/user"
	"github.com/artni96/GophKeeper/internal/service/login"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const DefaultError = "Internal Server Error"

// Handler represents the Login handler instance.
type Handler struct {
	pb.UnimplementedLoginServiceServer
	Service login.ServiceI
	Logger  *zap.Logger
}

// NewHandler initializes and returns the Login Handler instance.
func NewHandler(service login.ServiceI, logger *zap.Logger) *Handler {
	return &Handler{
		Service: service,
		Logger:  logger,
	}
}

// CreateLogin is a handler to create a new Login entity.
func (h *Handler) CreateLogin(ctx context.Context, req *pb.LoginCreateRequest) (*pb.LoginCreateResponse, error) {
	resp := &pb.LoginCreateResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityToCreate := loginmodel.CreateLoginRequest{
		Login:       req.GetLogin(),
		Password:    req.GetPassword(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		URL:         req.GetUrl(),
		UserID:      userID,
	}
	err := h.Service.Create(ctx, entityToCreate)
	if err != nil {
		return resp, status.Error(codes.Internal, DefaultError)
	}
	return resp, nil
}

// GetLogin is a handler to get a LoginEntity by its number and author.
func (h *Handler) GetLogin(ctx context.Context, req *pb.LoginGetRequest) (*pb.LoginGetResponse, error) {
	resp := &pb.LoginGetResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	loginNumber := req.GetNumber()
	dbEntity, err := h.Service.GetByNumber(ctx, loginNumber, userID)
	if err != nil {
		e := &userrepo.ErrUserNotFound
		if errors.As(err, e) {
			return resp, status.Error(codes.NotFound, err.Error())
		}
		return resp, status.Error(codes.Internal, DefaultError)
	}
	h.prepareResponse(*dbEntity, resp)
	return resp, nil
}

// prepareResponse sets pb.LoginGetResponse field values from a Login database entity - Login service result.
func (h *Handler) prepareResponse(dbEntity loginmodel.Login, resp *pb.LoginGetResponse) *pb.LoginGetResponse {
	resp.SetUrl(dbEntity.URL)
	resp.SetTitle(dbEntity.Title)
	resp.SetDescription(dbEntity.Description)
	resp.SetNumber(dbEntity.Number)
	resp.SetLogin(dbEntity.Login)
	resp.SetPassword(dbEntity.Password)
	resp.SetCreatedAt(dbEntity.CreatedAt.Format(time.RFC3339))
	if !dbEntity.UpdatedAt.IsZero() {
		resp.SetUpdatedAt(dbEntity.UpdatedAt.Format(time.RFC3339))
	}
	return resp
}
