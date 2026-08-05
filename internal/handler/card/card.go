package card

import (
	"context"
	"errors"
	"time"

	pb "github.com/artni96/GophKeeper/api/proto/cards"
	"github.com/artni96/GophKeeper/internal/constants"
	"github.com/artni96/GophKeeper/internal/interceptors"
	cardmodel "github.com/artni96/GophKeeper/internal/model/card"
	"github.com/artni96/GophKeeper/internal/service/card"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler represents the Card handler instance.
type Handler struct {
	pb.UnimplementedCardServiceServer
	Service card.ServiceI
	Logger  *zap.Logger
}

// NewHandler initializes and returns the Card Handler instance.
func NewHandler(service card.ServiceI, logger *zap.Logger) *Handler {
	return &Handler{
		Service: service,
		Logger:  logger,
	}
}

// CreateCard is a handler to create a new Card entity.
func (h *Handler) CreateCard(ctx context.Context, req *pb.CardCreateRequest) (*pb.CardCreateResponse, error) {
	resp := &pb.CardCreateResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityToCreate := cardmodel.CreateCardRequest{
		PAN:         req.GetPan(),
		Holder:      req.GetHolder(),
		ExpiryDate:  req.GetExpiryDate(),
		CVV:         req.GetCvv(),
		PIN:         req.GetPin(),
		Bank:        req.GetBank(),
		Brand:       req.GetBrand(),
		UserID:      userID,
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	}
	if err := entityToCreate.Validate(); err != nil {
		return resp, status.Error(codes.InvalidArgument, err.Error())
	}
	number, err := h.Service.Create(ctx, entityToCreate)
	if err != nil {
		return resp, status.Errorf(codes.Internal, "%v", err)
	}
	resp.SetNumber(number)
	return resp, nil
}

// GetCard is a handler to get a Card entity by its number and author.
func (h *Handler) GetCard(ctx context.Context, req *pb.CardGetRequest) (*pb.CardGetResponse, error) {
	resp := &pb.CardGetResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityNumber := req.GetNumber()
	dbEntity, err := h.Service.GetByNumber(ctx, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to get card record", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "card record with number %d not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}

	resp.SetPan(dbEntity.PAN)
	resp.SetHolder(dbEntity.Holder)
	resp.SetExpiryDate(dbEntity.ExpiryDate)
	resp.SetCvv(dbEntity.CVV)
	resp.SetPin(dbEntity.PIN)
	resp.SetBank(dbEntity.Bank)
	resp.SetBrand(dbEntity.Brand)
	resp.SetTitle(dbEntity.Title)
	resp.SetDescription(dbEntity.Description)
	resp.SetCreatedAt(dbEntity.CreatedAt.Format(time.RFC3339))
	if !dbEntity.UpdatedAt.IsZero() {
		resp.SetUpdatedAt(dbEntity.UpdatedAt.Format(time.RFC3339))
	}
	resp.SetNumber(dbEntity.Number)
	return resp, nil
}

// UpdateCard is a handler to update a Card entity by its number (only for authors).
func (h *Handler) UpdateCard(ctx context.Context, req *pb.CardUpdateRequest) (*pb.CardUpdateResponse, error) {
	resp := &pb.CardUpdateResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityNumber := req.GetNumber()
	if entityNumber == 0 {
		return resp, status.Error(codes.NotFound, "number is required")
	}

	dataToUpdate := cardmodel.UpdateCardRequest{
		PAN:         req.GetPan(),
		Holder:      req.GetHolder(),
		ExpiryDate:  req.GetExpiryDate(),
		CVV:         req.GetCvv(),
		PIN:         req.GetPin(),
		Bank:        req.GetBank(),
		Brand:       req.GetBrand(),
		Title:       req.GetTitle(),
		Description: req.GetDescription(),
	}
	if err := dataToUpdate.Validate(); err != nil {
		return resp, status.Error(codes.InvalidArgument, err.Error())
	}

	err := h.Service.Update(ctx, dataToUpdate, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to update card record", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "card record with %d number not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}
	return resp, nil
}

// DeleteCard is a handler to delete a Card entity by its number (only for authors).
func (h *Handler) DeleteCard(ctx context.Context, req *pb.CardDeleteRequest) (*pb.CardDeleteResponse, error) {
	resp := &pb.CardDeleteResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}
	entityNumber := req.GetNumber()

	err := h.Service.Delete(ctx, entityNumber, userID)
	if err != nil {
		h.Logger.Info("failed to delete card record", zap.Error(err))
		if errors.Is(err, constants.ErrEntityNotFound) {
			return resp, status.Errorf(codes.NotFound, "card record with %d number not found", entityNumber)
		}
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}
	return resp, nil
}

// GetListCard is a handler to retrieve user's Cart entities list.
func (h *Handler) GetListCard(ctx context.Context, req *pb.CardGetListRequest) (*pb.CardGetListResponse, error) {
	resp := &pb.CardGetListResponse{}
	userID, ok := interceptors.GetUserIDFromContext(ctx)
	if !ok {
		return resp, status.Errorf(codes.Unauthenticated, "not authenticated")
	}

	dbEntities, err := h.Service.GetList(ctx, userID)
	if err != nil {
		h.Logger.Info("failed to retrieve user's card records list", zap.Error(err), zap.String("user_id", userID.String()))
		return resp, status.Error(codes.Internal, constants.DefaultError)
	}

	pbList := make([]*pb.CardGetListItemResponse, 0)
	for _, entity := range dbEntities {
		i := &pb.CardGetListItemResponse{}
		i.SetNumber(entity.Number)
		i.SetTitle(entity.Title)
		i.SetDescription(entity.Description)
		pbList = append(pbList, i)
	}
	resp.SetCards(pbList)

	return resp, nil
}
