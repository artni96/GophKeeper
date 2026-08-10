package card

import (
	"context"
	"time"

	cardspb "github.com/artni96/GophKeeper/api/proto/cards"
	"github.com/artni96/GophKeeper/internal/client/config"
	"github.com/artni96/GophKeeper/internal/client/constants"
	"github.com/artni96/GophKeeper/internal/client/model/card"
	"github.com/artni96/GophKeeper/internal/client/model/common"
	cardrepo "github.com/artni96/GophKeeper/internal/client/repository/card"
	"github.com/artni96/GophKeeper/internal/client/utils"
	"google.golang.org/grpc"
)

// Service implements the Card client service to manage card-related data business logic.
type Service struct {
	cfg    *config.Config
	Client cardspb.CardServiceClient
	repo   cardrepo.RepositoryI
	state  *config.State
}

// NewService initializes and returns the new Card service instance.
func NewService(cfg *config.Config, conn *grpc.ClientConn, repo cardrepo.RepositoryI, state *config.State) *Service {
	return &Service{
		cfg:    cfg,
		Client: cardspb.NewCardServiceClient(conn),
		repo:   repo,
		state:  state,
	}
}

// Add adds a new Card entity in the repository.
func (s *Service) Add(ctx context.Context, entityNumber uint64) error {
	cardGetReq := &cardspb.CardGetRequest{}
	cardGetReq.SetNumber(entityNumber)

	pbEntity, err := s.Client.GetCard(ctx, cardGetReq)
	if err != nil {
		return err
	}

	createdAt, err := time.Parse(time.RFC3339, pbEntity.GetCreatedAt())
	if err != nil {
		return err
	}

	var updatedAt time.Time
	updatedAt, err = time.Parse(time.RFC3339, pbEntity.GetUpdatedAt())
	if err != nil {
		updatedAt = time.Time{}
	}

	panNonce := pbEntity.GetPanNonce()
	panKeyID := pbEntity.GetPanKeyId()
	panAesKey := s.state.Keys[panKeyID]

	decryptedPAN, err := utils.DecryptField(pbEntity.GetPan(), panAesKey, panNonce)
	if err != nil {
		return err
	}

	holderNonce := pbEntity.GetHolderNonce()
	holderKeyID := pbEntity.GetHolderKeyId()
	holderAesKey := s.state.Keys[holderKeyID]

	decryptedHolder, err := utils.DecryptField(pbEntity.GetHolder(), holderAesKey, holderNonce)
	if err != nil {
		return err
	}

	expiryDateNonce := pbEntity.GetExpiryDateNonce()
	expiryDateKeyID := pbEntity.GetExpiryDateKeyId()
	expiryDateAesKey := s.state.Keys[expiryDateKeyID]

	decryptedExpiryDate, err := utils.DecryptField(pbEntity.GetExpiryDate(), expiryDateAesKey, expiryDateNonce)
	if err != nil {
		return err
	}

	cvvNonce := pbEntity.GetCvvNonce()
	cvvKeyID := pbEntity.GetCvvKeyId()
	cvvAesKey := s.state.Keys[cvvKeyID]

	decryptedCVV, err := utils.DecryptField(pbEntity.GetCvv(), cvvAesKey, cvvNonce)
	if err != nil {
		return err
	}

	pinNonce := pbEntity.GetPinNonce()
	pinKeyID := pbEntity.GetPinKeyId()
	pinAesKey := s.state.Keys[pinKeyID]

	decryptedPIN, err := utils.DecryptField(pbEntity.GetPin(), pinAesKey, pinNonce)
	if err != nil {
		return err
	}

	entity := card.Card{
		PAN:         string(decryptedPAN),
		Holder:      string(decryptedHolder),
		ExpiryDate:  string(decryptedExpiryDate),
		CVV:         string(decryptedCVV),
		PIN:         string(decryptedPIN),
		Bank:        pbEntity.GetBank(),
		Brand:       pbEntity.GetBrand(),
		Title:       pbEntity.GetTitle(),
		Number:      pbEntity.GetNumber(),
		Description: pbEntity.GetDescription(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}
	s.repo.Add(entity)
	return nil
}

// AddBatch adds a list of new Card entities in the repository.
func (s *Service) AddBatch(ctx context.Context) error {
	listReq := &cardspb.CardGetListRequest{}

	listResp, err := s.Client.GetListCard(ctx, listReq)
	if err != nil {
		return err
	}

	pbEntities := listResp.GetCards()

	var entities []card.Card
	for _, pbEntity := range pbEntities {
		createdAt, err := time.Parse(time.RFC3339, pbEntity.GetCreatedAt())
		if err != nil {
			return err
		}

		var updatedAt time.Time
		updatedAt, err = time.Parse(time.RFC3339, pbEntity.GetUpdatedAt())
		if err != nil {
			updatedAt = time.Time{}
		}

		panNonce := pbEntity.GetPanNonce()
		panKeyID := pbEntity.GetPanKeyId()
		panAesKey := s.state.Keys[panKeyID]

		decryptedPAN, err := utils.DecryptField(pbEntity.GetPan(), panAesKey, panNonce)
		if err != nil {
			return err
		}

		holderNonce := pbEntity.GetHolderNonce()
		holderKeyID := pbEntity.GetHolderKeyId()
		holderAesKey := s.state.Keys[holderKeyID]

		decryptedHolder, err := utils.DecryptField(pbEntity.GetHolder(), holderAesKey, holderNonce)
		if err != nil {
			return err
		}

		expiryDateNonce := pbEntity.GetExpiryDateNonce()
		expiryDateKeyID := pbEntity.GetExpiryDateKeyId()
		expiryDateAesKey := s.state.Keys[expiryDateKeyID]

		decryptedExpiryDate, err := utils.DecryptField(pbEntity.GetExpiryDate(), expiryDateAesKey, expiryDateNonce)
		if err != nil {
			return err
		}

		cvvNonce := pbEntity.GetCvvNonce()
		cvvKeyID := pbEntity.GetCvvKeyId()
		cvvAesKey := s.state.Keys[cvvKeyID]

		decryptedCVV, err := utils.DecryptField(pbEntity.GetCvv(), cvvAesKey, cvvNonce)
		if err != nil {
			return err
		}

		pinNonce := pbEntity.GetPinNonce()
		pinKeyID := pbEntity.GetPinKeyId()
		pinAesKey := s.state.Keys[pinKeyID]

		decryptedPIN, err := utils.DecryptField(pbEntity.GetPin(), pinAesKey, pinNonce)
		if err != nil {
			return err
		}

		entity := card.Card{
			PAN:         string(decryptedPAN),
			Holder:      string(decryptedHolder),
			ExpiryDate:  string(decryptedExpiryDate),
			CVV:         string(decryptedCVV),
			PIN:         string(decryptedPIN),
			Bank:        pbEntity.GetBank(),
			Brand:       pbEntity.GetBrand(),
			Title:       pbEntity.GetTitle(),
			Number:      pbEntity.GetNumber(),
			Description: pbEntity.GetDescription(),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
		}
		entities = append(entities, entity)

	}
	s.repo.AddBatch(entities)
	return nil
}

// Get returns the Card entity by its number from the repository.
func (s *Service) Get(entityNumber uint64) (card.Card, error) {
	entity, err := s.repo.Get(entityNumber)
	if err != nil {
		return card.Card{}, err
	}
	return entity, nil
}

// GetList returns the list of Card entities from the repository.
func (s *Service) GetList() []common.Entity {
	return s.repo.GetList()
}

// Update updates the Card entity in the repository.
func (s *Service) Update(ctx context.Context, entityNumber uint64) error {
	req := &cardspb.CardGetRequest{}
	req.SetNumber(entityNumber)
	pbEntity, err := s.Client.GetCard(ctx, req)
	if err != nil {
		return constants.ErrEntityNotFound
	}

	createdAt, err := time.Parse(time.RFC3339, pbEntity.GetCreatedAt())
	if err != nil {
		return err
	}
	updatedAt, err := time.Parse(time.RFC3339, pbEntity.GetUpdatedAt())
	if err != nil {
		return err
	}

	panNonce := pbEntity.GetPanNonce()
	panKeyID := pbEntity.GetPanKeyId()
	panAesKey := s.state.Keys[panKeyID]

	decryptedPAN, err := utils.DecryptField(pbEntity.GetPan(), panAesKey, panNonce)
	if err != nil {
		return err
	}

	holderNonce := pbEntity.GetHolderNonce()
	holderKeyID := pbEntity.GetHolderKeyId()
	holderAesKey := s.state.Keys[holderKeyID]

	decryptedHolder, err := utils.DecryptField(pbEntity.GetHolder(), holderAesKey, holderNonce)
	if err != nil {
		return err
	}

	expiryDateNonce := pbEntity.GetExpiryDateNonce()
	expiryDateKeyID := pbEntity.GetExpiryDateKeyId()
	expiryDateAesKey := s.state.Keys[expiryDateKeyID]

	decryptedExpiryDate, err := utils.DecryptField(pbEntity.GetExpiryDate(), expiryDateAesKey, expiryDateNonce)
	if err != nil {
		return err
	}

	cvvNonce := pbEntity.GetCvvNonce()
	cvvKeyID := pbEntity.GetCvvKeyId()
	cvvAesKey := s.state.Keys[cvvKeyID]

	decryptedCVV, err := utils.DecryptField(pbEntity.GetCvv(), cvvAesKey, cvvNonce)
	if err != nil {
		return err
	}

	pinNonce := pbEntity.GetPinNonce()
	pinKeyID := pbEntity.GetPinKeyId()
	pinAesKey := s.state.Keys[pinKeyID]

	decryptedPIN, err := utils.DecryptField(pbEntity.GetPin(), pinAesKey, pinNonce)
	if err != nil {
		return err
	}

	entity := card.Card{
		PAN:         string(decryptedPAN),
		Holder:      string(decryptedHolder),
		ExpiryDate:  string(decryptedExpiryDate),
		CVV:         string(decryptedCVV),
		PIN:         string(decryptedPIN),
		Bank:        pbEntity.GetBank(),
		Brand:       pbEntity.GetBrand(),
		Title:       pbEntity.GetTitle(),
		Description: pbEntity.GetDescription(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Number:      entityNumber,
	}
	err = s.repo.Update(entity)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes the Card entity by its number from the repository.
func (s *Service) Delete(entityNumber uint64) error {
	err := s.repo.Delete(entityNumber)
	if err != nil {
		return err
	}
	return nil
}
