package card

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CreateCardRequest struct {
	PAN      *wrapperspb.BytesValue
	PANNonce *wrapperspb.BytesValue
	PANKeyID *wrapperspb.UInt64Value

	Holder      *wrapperspb.BytesValue
	HolderNonce *wrapperspb.BytesValue
	HolderKeyID *wrapperspb.UInt64Value

	ExpiryDate      *wrapperspb.BytesValue
	ExpiryDateNonce *wrapperspb.BytesValue
	ExpiryDateKeyID *wrapperspb.UInt64Value

	CVV      *wrapperspb.BytesValue
	CVVNonce *wrapperspb.BytesValue
	CVVKeyID *wrapperspb.UInt64Value

	PIN      *wrapperspb.BytesValue
	PINNonce *wrapperspb.BytesValue
	PINKeyID *wrapperspb.UInt64Value

	Bank        *wrapperspb.StringValue
	Brand       *wrapperspb.StringValue
	UserID      uuid.UUID
	Title       *wrapperspb.StringValue
	Description *wrapperspb.StringValue
}

type CreateCard struct {
	PAN      []byte `gorm:"column:hashed_pan"`
	PanNonce []byte
	PanKeyID uint64

	Holder      []byte `gorm:"column:hashed_holder"`
	HolderNonce []byte
	HolderKeyID uint64

	ExpiryDate      []byte `gorm:"column:hashed_expiry_date"`
	ExpiryDateNonce []byte
	ExpiryDateKeyID uint64

	CVV      []byte `gorm:"column:hashed_cvv"`
	CVVNonce []byte
	CVVKeyID uint64

	PIN      []byte `gorm:"column:hashed_pin"`
	PINNonce []byte
	PINKeyID uint64

	Bank        string
	Brand       string
	UserID      uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time
	Number      uint64
}

type Card struct {
	PAN      []byte `gorm:"column:hashed_pan"`
	PanNonce []byte
	PanKeyID uint64

	Holder      []byte `gorm:"column:hashed_holder"`
	HolderNonce []byte
	HolderKeyID uint64

	ExpiryDate      []byte `gorm:"column:hashed_expiry_date"`
	ExpiryDateNonce []byte
	ExpiryDateKeyID uint64

	CVV      []byte `gorm:"column:hashed_cvv"`
	CVVNonce []byte
	CVVKeyID uint64

	PIN         []byte `gorm:"column:hashed_pin"`
	PINNonce    []byte
	PINKeyID    uint64
	Bank        string
	Brand       string
	UserID      uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Number      uint64
}

type UpdateCardRequest struct {
	PAN      *wrapperspb.BytesValue
	PANNonce *wrapperspb.BytesValue
	PANKeyID *wrapperspb.UInt64Value

	Holder      *wrapperspb.BytesValue
	HolderNonce *wrapperspb.BytesValue
	HolderKeyID *wrapperspb.UInt64Value

	ExpiryDate      *wrapperspb.BytesValue
	ExpiryDateNonce *wrapperspb.BytesValue
	ExpiryDateKeyID *wrapperspb.UInt64Value

	CVV      *wrapperspb.BytesValue
	CVVNonce *wrapperspb.BytesValue
	CVVKeyID *wrapperspb.UInt64Value

	PIN      *wrapperspb.BytesValue
	PINNonce *wrapperspb.BytesValue
	PINKeyID *wrapperspb.UInt64Value

	Bank        *wrapperspb.StringValue
	Brand       *wrapperspb.StringValue
	Title       *wrapperspb.StringValue
	Description *wrapperspb.StringValue
}

type UpdateCard struct {
	PAN      []byte `gorm:"column:hashed_pan"`
	PanNonce []byte
	PanKeyID uint64

	Holder      []byte `gorm:"column:hashed_holder"`
	HolderNonce []byte
	HolderKeyID uint64

	ExpiryDate      []byte `gorm:"column:hashed_expiry_date"`
	ExpiryDateNonce []byte
	ExpiryDateKeyID uint64

	CVV      []byte `gorm:"column:hashed_cvv"`
	CVVNonce []byte
	CVVKeyID uint64

	PIN      []byte `gorm:"column:hashed_pin"`
	PINNonce []byte
	PINKeyID uint64

	Bank        string
	Brand       string
	Title       string
	Description string
	UpdatedAt   time.Time
}
