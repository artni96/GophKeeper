package login

import (
	"context"
	"errors"
	"sync"
	"time"

	pb "github.com/artni96/GophKeeper/api/proto/logins"
	userspb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/server/constants"
	"github.com/artni96/GophKeeper/internal/server/interceptors"
	loginmodel "github.com/artni96/GophKeeper/internal/server/model/login"
	"github.com/artni96/GophKeeper/internal/server/service/login"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Handler represents the Login handler instance.
type Handler struct {
	pb.UnimplementedLoginServiceServer
	Service login.ServiceI
	Logger  *zap.Logger
	streams map[uuid.UUID][]chan *userspb.UpdateNotification
	mu      sync.Mutex
}

// NewHandler initializes and returns the Login Handler instance.
func NewHandler(service login.ServiceI, logger *zap.Logger, streams map[uuid.UUID][]chan *userspb.UpdateNotification) *Handler {
	return &Handler{
		Service: service,
		Logger:  logger,
		streams: streams,
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
	if err := entityToCreate.Validate(); err != nil {
		return resp, status.Error(codes.InvalidArgument, err.Error())
	}
	entityNumber, err := h.Service.Create(ctx, entityToCreate)
	if err != nil {
		if errors.Is(err, constants.ErrEntityAlreadyExists) {
			return resp, status.Error(codes.AlreadyExists, err.Error())
		}
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}
	resp.SetNumber(entityNumber)

	h.mu.Lock()
	userStreams := h.streams[userID]
	h.mu.Unlock()

	notification := &userspb.UpdateNotification{}
	notification.SetUpdatedAt(timestamppb.Now())
	notification.SetNumber(entityNumber)
	notification.SetEntityType(2)
	notification.SetActionType(2)
	for _, stream := range userStreams {
		select {
		case stream <- notification:
			h.Logger.Debug("update notification added to stream channel")
		default:
			h.Logger.Debug("stream channel is full")
		}
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
		h.Logger.Info("failed to the get login record", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "login record with number %d not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}

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
	return resp, nil
}

// UpdateLogin is a handler to update a Login entity by its number (only for authors).
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
	if err := dataToUpdate.Validate(); err != nil {
		return resp, status.Error(codes.InvalidArgument, err.Error())
	}

	err := h.Service.Update(ctx, dataToUpdate, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to update login record", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "login record with %d number not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}

	h.mu.Lock()
	userStreams := h.streams[userID]
	h.mu.Unlock()

	notification := &userspb.UpdateNotification{}
	notification.SetUpdatedAt(timestamppb.Now())
	notification.SetNumber(entityNumber)
	notification.SetEntityType(2)
	notification.SetActionType(2)
	for _, stream := range userStreams {
		select {
		case stream <- notification:
			h.Logger.Debug("update notification added to stream channel")
		default:
			h.Logger.Debug("stream channel is full")
		}
	}

	return resp, nil
}

// DeleteLogin is a handler to delete a Login entity by its number (only for authors).
func (h *Handler) DeleteLogin(ctx context.Context, req *pb.LoginDeleteRequest) (*pb.LoginDeleteResponse, error) {
	resp := &pb.LoginDeleteResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityNumber := req.GetNumber()

	err := h.Service.Delete(ctx, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to delete login record", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "login record with %d number not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}

	h.mu.Lock()
	userStreams := h.streams[userID]
	h.mu.Unlock()

	notification := &userspb.UpdateNotification{}
	notification.SetNumber(entityNumber)
	notification.SetEntityType(2)
	notification.SetActionType(1)
	for _, stream := range userStreams {
		select {
		case stream <- notification:
			h.Logger.Debug("update notification added to stream channel")
		default:
			h.Logger.Debug("stream channel is full")
		}
	}

	return resp, nil
}

// GetListLogin is a handler to retrieve user's Login entities list.
func (h *Handler) GetListLogin(ctx context.Context, req *pb.LoginGetListRequest) (*pb.LoginGetListResponse, error) {
	resp := &pb.LoginGetListResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}

	dbEntities, err := h.Service.GetList(ctx, userID)
	if err != nil {
		h.Logger.Info("failed to retrieve user's login records list", zap.Error(err), zap.String("user_id", userID.String()))
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}

	pbList := make([]*pb.LoginGetResponse, 0)
	for _, entity := range dbEntities {
		i := &pb.LoginGetResponse{}
		i.SetTitle(entity.Title)
		i.SetNumber(entity.Number)
		i.SetUrl(entity.URL)
		i.SetDescription(entity.Description)
		i.SetLogin(entity.Login)
		i.SetPassword(entity.Password)
		i.SetCreatedAt(entity.CreatedAt.Format(time.RFC3339))
		if !entity.UpdatedAt.IsZero() {
			i.SetUpdatedAt(entity.UpdatedAt.Format(time.RFC3339))
		}
		pbList = append(pbList, i)
	}
	resp.SetLogins(pbList)

	return resp, nil
}
