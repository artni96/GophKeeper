package login

import (
	"fmt"
	"time"

	"github.com/artni96/GophKeeper/internal/server/model/common"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// CreateLoginRequest is a model-mediator for transferring the gRPC request data to the Login service.
type CreateLoginRequest struct {
	Login    common.PBEncryptedField
	Password common.PBEncryptedField

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
	UserID   uuid.UUID
	Login    common.EncryptedField `gorm:"embedded;embeddedPrefix:login_"`
	Password common.EncryptedField `gorm:"embedded;embeddedPrefix:password_"`

	Title       string
	Number      uint64
	URL         string
	Description string
	CreatedAt   time.Time
}

type Login struct {
	ID       uuid.UUID
	UserID   uuid.UUID
	Login    common.EncryptedField `gorm:"embedded;embeddedPrefix:login_"`
	Password common.EncryptedField `gorm:"embedded;embeddedPrefix:password_"`

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

	Login    common.PBEncryptedField
	Password common.PBEncryptedField
}

func (u *UpdateLoginRequest) Validate() error {
	if u.Title != nil && u.Title.String() == "" {
		return fmt.Errorf("title field is required")
	}
	return nil
}

type UpdateLogin struct {
	Login    common.EncryptedField `gorm:"embedded;embeddedPrefix:login_"`
	Password common.EncryptedField `gorm:"embedded;embeddedPrefix:password_"`

	Title       string
	URL         string
	Description string
	UpdatedAt   time.Time
}
