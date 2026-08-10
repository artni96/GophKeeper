package text

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CreateTextRequest struct {
	Title       *wrapperspb.StringValue
	Description *wrapperspb.StringValue

	Text  *wrapperspb.BytesValue
	Nonce *wrapperspb.BytesValue
	KeyID *wrapperspb.UInt64Value

	UserID uuid.UUID
}

type CreateText struct {
	Title       string
	Description string

	Text  []byte `gorm:"column:hashed_text"`
	Nonce []byte
	KeyID uint64

	UserID    uuid.UUID
	CreatedAt time.Time
	Number    uint64
}

type Text struct {
	Title       string
	Description string

	Text  []byte `gorm:"column:hashed_text"`
	Nonce []byte
	KeyID uint64

	UserID    uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Number    uint64
}

type UpdateTextRequest struct {
	Title       *wrapperspb.StringValue
	Description *wrapperspb.StringValue

	Text  *wrapperspb.BytesValue
	Nonce *wrapperspb.BytesValue
	KeyID *wrapperspb.UInt64Value
}

type UpdateText struct {
	Title       string
	Description string

	Text  []byte `gorm:"column:hashed_text"`
	Nonce []byte
	KeyID uint64

	UpdatedAt time.Time
}
