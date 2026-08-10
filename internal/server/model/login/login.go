package login

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// CreateLoginRequest is a model-mediator for transferring the gRPC request data to the Login service.
type CreateLoginRequest struct {
	Login      *wrapperspb.BytesValue
	LoginNonce *wrapperspb.BytesValue
	LoginKeyID *wrapperspb.UInt64Value

	Password      *wrapperspb.BytesValue
	PasswordNonce *wrapperspb.BytesValue
	PasswordKeyID *wrapperspb.UInt64Value

	UserID      uuid.UUID
	Title       *wrapperspb.StringValue
	URL         *wrapperspb.StringValue
	Description *wrapperspb.StringValue
}

func (m *CreateLoginRequest) Validate() error {
	if m.Title.Value == "" {
		return fmt.Errorf("title field is required")
	}
	return nil
}

// CreateLogin is a model-mediator for transferring the Login creation data to the Login repository.
type CreateLogin struct {
	UserID     uuid.UUID
	Login      []byte `gorm:"column:hashed_login;not null"`
	LoginNonce []byte
	LoginKeyID uint64

	Password      []byte `gorm:"column:hashed_password;not null"`
	PasswordNonce []byte
	PasswordKeyID uint64
	Title         string
	Number        uint64
	URL           string
	Description   string
	CreatedAt     time.Time
}

type Login struct {
	ID         uuid.UUID
	UserID     uuid.UUID
	Login      []byte `gorm:"column:hashed_login"`
	LoginNonce []byte
	LoginKeyID uint64

	Password      []byte `gorm:"column:hashed_password"`
	PasswordNonce []byte
	PasswordKeyID uint64

	Title       string
	Number      uint64
	URL         string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type UpdateLoginRequest struct {
	Title       *wrapperspb.StringValue
	URL         *wrapperspb.StringValue
	Description *wrapperspb.StringValue

	Login      *wrapperspb.BytesValue
	LoginNonce *wrapperspb.BytesValue
	LoginKeyID *wrapperspb.UInt64Value

	Password      *wrapperspb.BytesValue
	PasswordNonce *wrapperspb.BytesValue
	PasswordKeyID *wrapperspb.UInt64Value
}

func (u *UpdateLoginRequest) Validate() error {
	if u.Title != nil && u.Title.String() == "" {
		return fmt.Errorf("title field is required")
	}
	return nil
}

type UpdateLogin struct {
	Login      []byte `gorm:"column:hashed_login"`
	LoginNonce []byte
	LoginKeyID uint64

	Password      []byte `gorm:"column:hashed_password"`
	PasswordNonce []byte
	PasswordKeyID uint64

	Title       string
	URL         string
	Description string
	UpdatedAt   time.Time
}

type GetListLoginResponse struct {
	Title       string
	Description string
	Number      uint64
}
