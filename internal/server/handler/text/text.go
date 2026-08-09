package text

import (
	"context"
	"errors"
	"sync"
	"time"

	pb "github.com/artni96/GophKeeper/api/proto/texts"
	userspb "github.com/artni96/GophKeeper/api/proto/users"
	"github.com/artni96/GophKeeper/internal/server/constants"
	"github.com/artni96/GophKeeper/internal/server/interceptors"
	textmodel "github.com/artni96/GophKeeper/internal/server/model/text"
	"github.com/artni96/GophKeeper/internal/server/service/text"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Handler represents the Text handler instance.
type Handler struct {
	pb.UnimplementedTextServiceServer
	Service text.ServiceI
	Logger  *zap.Logger
	streams map[uuid.UUID][]chan *userspb.UpdateNotification
	mu      sync.RWMutex
}

// NewHandler initializes and returns the Text Handler instance.
func NewHandler(service text.ServiceI, logger *zap.Logger, streams map[uuid.UUID][]chan *userspb.UpdateNotification) *Handler {
	return &Handler{
		Service: service,
		Logger:  logger,
		streams: streams,
	}
}

// CreateText is a handler to create a new Text entity.
func (h *Handler) CreateText(ctx context.Context, req *pb.TextCreateRequest) (*pb.TextCreateResponse, error) {
	resp := &pb.TextCreateResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityToCreate := textmodel.CreateTextRequest{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Text:        req.GetText(),
		UserID:      userID,
	}
	entityNumber, err := h.Service.Create(ctx, entityToCreate)
	if err != nil {
		if errors.Is(err, constants.ErrEntityAlreadyExists) {
			return resp, status.Errorf(codes.AlreadyExists, "%v", err)
		}
		return resp, status.Errorf(codes.Internal, "%v", err)

	}
	resp.SetNumber(entityNumber)

	h.mu.Lock()
	userStreams := h.streams[userID]
	h.mu.Unlock()

	notification := &userspb.UpdateNotification{}
	notification.SetUpdatedAt(timestamppb.Now())
	notification.SetNumber(entityNumber)
	notification.SetEntityType(userspb.EntityType_Text)
	notification.SetActionType(userspb.ActionType_Create)
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

// GetText is a handler to get a Text entity by its number and author.
func (h *Handler) GetText(ctx context.Context, req *pb.TextGetRequest) (*pb.TextGetResponse, error) {
	resp := &pb.TextGetResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityNumber := req.GetNumber()
	dbEntity, err := h.Service.GetByNumber(ctx, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to get text record", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "text record with number %d not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}

	resp.SetTitle(dbEntity.Title)
	resp.SetDescription(dbEntity.Description)
	resp.SetText(dbEntity.Text)
	resp.SetCreatedAt(dbEntity.CreatedAt.Format(time.RFC3339))
	if !dbEntity.UpdatedAt.IsZero() {
		resp.SetUpdatedAt(dbEntity.UpdatedAt.Format(time.RFC3339))
	}
	resp.SetNumber(dbEntity.Number)
	return resp, nil
}

// UpdateText is a handler to update a Text entity by its number (only for authors).
func (h *Handler) UpdateText(ctx context.Context, req *pb.TextUpdateRequest) (*pb.TextUpdateResponse, error) {
	resp := &pb.TextUpdateResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityNumber := req.GetNumber()
	if entityNumber == 0 {
		return resp, status.Error(codes.NotFound, "number is required")
	}

	dataToUpdate := textmodel.UpdateTextRequest{
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
		Text:        req.GetText(),
	}

	err := h.Service.Update(ctx, dataToUpdate, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to update text record", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "text record with number %d not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}

	h.mu.Lock()
	userStreams := h.streams[userID]
	h.mu.Unlock()

	notification := &userspb.UpdateNotification{}
	notification.SetUpdatedAt(timestamppb.Now())
	notification.SetNumber(entityNumber)
	notification.SetEntityType(userspb.EntityType_Text)
	notification.SetActionType(userspb.ActionType_Update)
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

// DeleteText is a handler to delete a Text entity by its number (only for authors).
func (h *Handler) DeleteText(ctx context.Context, req *pb.TextDeleteRequest) (*pb.TextDeleteResponse, error) {
	resp := &pb.TextDeleteResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityNumber := req.GetNumber()

	err := h.Service.Delete(ctx, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to delete text record", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "text record with %d number not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}

	h.mu.Lock()
	userStreams := h.streams[userID]
	h.mu.Unlock()

	notification := &userspb.UpdateNotification{}
	notification.SetUpdatedAt(timestamppb.Now())
	notification.SetNumber(entityNumber)
	notification.SetEntityType(userspb.EntityType_Text)
	notification.SetActionType(userspb.ActionType_Delete)
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

// GetListText is a handler to retrieve user's Text entities list.
func (h *Handler) GetListText(ctx context.Context, req *pb.TextGetListRequest) (*pb.TextGetListResponse, error) {
	resp := &pb.TextGetListResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}

	dbEntities, err := h.Service.GetList(ctx, userID)
	if err != nil {
		h.Logger.Info("failed to retrieve user's text records list", zap.Error(err), zap.String("user_id", userID.String()))
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}

	pbList := make([]*pb.TextGetResponse, 0)
	for _, entity := range dbEntities {
		i := &pb.TextGetResponse{}
		i.SetTitle(entity.Title)
		i.SetDescription(entity.Description)
		i.SetText(entity.Text)
		i.SetNumber(entity.Number)
		i.SetCreatedAt(entity.CreatedAt.Format(time.RFC3339))
		if !entity.UpdatedAt.IsZero() {
			i.SetUpdatedAt(entity.UpdatedAt.Format(time.RFC3339))
		}
		pbList = append(pbList, i)
	}
	resp.SetTexts(pbList)
	return resp, nil
}
