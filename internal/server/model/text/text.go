package text

import (
	"time"

	"github.com/artni96/GophKeeper/internal/server/model/common"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CreateTextRequest struct {
	Title       *wrapperspb.StringValue
	Description *wrapperspb.StringValue

	Text common.PBEncryptedField

	UserID uuid.UUID
}

type CreateText struct {
	Title       string
	Description string

	Text common.EncryptedField `gorm:"embedded;embeddedPrefix:text_"`

	UserID    uuid.UUID
	CreatedAt time.Time
	Number    uint64
}

func (c *CreateText) GetUserID() uuid.UUID {
	return c.UserID
}

func (c *CreateText) SetNumber(number uint64) {
	c.Number = number
}

type Text struct {
	Title       string
	Description string

	Text common.EncryptedField `gorm:"embedded;embeddedPrefix:text_"`

	ID        uuid.UUID
	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Number    uint64
}

type UpdateTextRequest struct {
	Title       *wrapperspb.StringValue
	Description *wrapperspb.StringValue

	Text common.PBEncryptedField
}

type UpdateText struct {
	Title       string
	Description string

	Text common.EncryptedField `gorm:"embedded;embeddedPrefix:text_"`

	UpdatedAt time.Time
}
