package login

import (
	"context"
	"errors"
	"time"

	pb "github.com/artni96/GophKeeper/api/proto/logins"
	"github.com/artni96/GophKeeper/internal/constants"
	"github.com/artni96/GophKeeper/internal/interceptors"
	loginmodel "github.com/artni96/GophKeeper/internal/model/login"
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
		if errors.Is(err, constants.ErrInvalidRequest) {
			return resp, status.Error(codes.InvalidArgument, err.Error())
		} else if errors.Is(err, constants.ErrEntityAlreadyExists) {
			return resp, status.Error(codes.AlreadyExists, err.Error())
		}
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
	entityNumber := req.GetNumber()
	dbEntity, err := h.Service.GetByNumber(ctx, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to get login", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "login with number %d not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, DefaultError)
	}
	h.prepareResponse(*dbEntity, resp)
	return resp, nil
}

// UpdateLogin is a handler to update the Login entity by its number (only for authors).
func (h *Handler) UpdateLogin(ctx context.Context, req *pb.LoginUpdateRequest) (*pb.LoginUpdateResponse, error) {
	resp := &pb.LoginUpdateResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityNumber := req.GetNumber()

	dataToUpdate := loginmodel.UpdateLoginRequest{
		Login:       req.GetLogin(),
		Password:    req.GetPassword(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		URL:         req.GetUrl(),
	}

	err := h.Service.Update(ctx, dataToUpdate, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to update login", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "login with %d number not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, DefaultError)
	}
	return resp, nil
}

// DeleteLogin is a handler to delete the Login entity by its number (only for authors).
func (h *Handler) DeleteLogin(ctx context.Context, req *pb.LoginDeleteRequest) (*pb.LoginDeleteResponse, error) {
	resp := &pb.LoginDeleteResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityNumber := req.GetNumber()

	err := h.Service.Delete(ctx, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to delete login", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "login with %d number not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, DefaultError)
	}
	return resp, nil
}

func (h *Handler) GetListLogin(ctx context.Context, req *pb.LoginGetListRequest) (*pb.LoginGetListResponse, error) {
	resp := &pb.LoginGetListResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}

	dbEntities, err := h.Service.GetList(ctx, userID)
	if err != nil {
		h.Logger.Info("failed to retrieve user's login list", zap.Error(err), zap.String("user_id", userID.String()))
		return resp, status.Error(codes.Internal, DefaultError)
	}

	pbList := make([]*pb.LoginGetListItemResponse, 0)
	for _, entity := range dbEntities {
		i := &pb.LoginGetListItemResponse{}
		i.SetNumber(entity.Number)
		i.SetTitle(entity.Title)
		i.SetDescription(entity.Description)
		pbList = append(pbList, i)
	}
	resp.SetLogins(pbList)

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
