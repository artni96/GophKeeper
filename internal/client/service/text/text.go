package text

import (
	"context"
	"time"

	textspb "github.com/artni96/GophKeeper/api/proto/texts"
	"github.com/artni96/GophKeeper/internal/client/config"
	"github.com/artni96/GophKeeper/internal/client/constants"
	"github.com/artni96/GophKeeper/internal/client/model/common"
	"github.com/artni96/GophKeeper/internal/client/model/text"
	textrepo "github.com/artni96/GophKeeper/internal/client/repository/text"
	"github.com/artni96/GophKeeper/internal/client/utils"
	"google.golang.org/grpc"
)

// Service implements the Text client service to manage text-related data business logic.
type Service struct {
	cfg    *config.Config
	Client textspb.TextServiceClient
	repo   textrepo.RepositoryI
	state  *config.State
}

// NewService initializes and returns the new Text service instance.
func NewService(cfg *config.Config, conn *grpc.ClientConn, repo textrepo.RepositoryI, state *config.State) *Service {
	return &Service{
		cfg:    cfg,
		Client: textspb.NewTextServiceClient(conn),
		repo:   repo,
		state:  state,
	}
}

// Add adds a new Text entity in the repository.
func (s *Service) Add(ctx context.Context, entityNumber uint64) error {
	textGetReq := &textspb.TextGetRequest{}
	textGetReq.SetNumber(entityNumber)

	pbEntity, err := s.Client.GetText(ctx, textGetReq)
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

	nonce := pbEntity.GetNonce()
	keyID := pbEntity.GetKeyId()
	aesKey := s.state.Keys[keyID]

	decryptedText, err := utils.DecryptField(pbEntity.GetText(), aesKey, nonce)
	if err != nil {
		return err
	}

	entity := text.Text{
		Title:       pbEntity.GetTitle(),
		Number:      pbEntity.GetNumber(),
		Description: pbEntity.GetDescription(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Text:        string(decryptedText),
	}
	s.repo.Add(entity)
	return nil
}

// AddBatch adds a list of new Text entities in the repository.
func (s *Service) AddBatch(ctx context.Context) error {
	listReq := &textspb.TextGetListRequest{}

	listResp, err := s.Client.GetListText(ctx, listReq)
	if err != nil {
		return err
	}

	pbEntities := listResp.GetTexts()

	var entities []text.Text
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

		nonce := pbEntity.GetNonce()
		keyID := pbEntity.GetKeyId()
		aesKey := s.state.Keys[keyID]

		decryptedText, err := utils.DecryptField(pbEntity.GetText(), aesKey, nonce)
		if err != nil {
			return err
		}

		entity := text.Text{
			Title:       pbEntity.GetTitle(),
			Number:      pbEntity.GetNumber(),
			Description: pbEntity.GetDescription(),
			CreatedAt:   createdAt,
			UpdatedAt:   updatedAt,
			Text:        string(decryptedText),
		}
		entities = append(entities, entity)

	}
	s.repo.AddBatch(entities)
	return nil
}

// Get returns the Text entity by its number in the repository.
func (s *Service) Get(entityNumber uint64) (text.Text, error) {
	entity, err := s.repo.Get(entityNumber)
	if err != nil {
		return text.Text{}, err
	}
	return entity, nil
}

// GetList returns the list of Text entities from the repository.
func (s *Service) GetList() []common.Entity {
	return s.repo.GetList()
}

// Update updates the Text entity in the repository.
func (s *Service) Update(ctx context.Context, entityNumber uint64) error {
	req := &textspb.TextGetRequest{}
	req.SetNumber(entityNumber)
	pbEntity, err := s.Client.GetText(ctx, req)
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

	nonce := pbEntity.GetNonce()
	keyID := pbEntity.GetKeyId()
	aesKey := s.state.Keys[keyID]

	decryptedText, err := utils.DecryptField(pbEntity.GetText(), aesKey, nonce)
	if err != nil {
		return err
	}

	entity := text.Text{
		Title:       pbEntity.GetTitle(),
		Number:      pbEntity.GetNumber(),
		Description: pbEntity.GetDescription(),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Text:        string(decryptedText),
	}
	err = s.repo.Update(entity)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes the Text entity by its number from the repository.
func (s *Service) Delete(entityNumber uint64) error {
	err := s.repo.Delete(entityNumber)
	if err != nil {
		return err
	}
	return nil
}
