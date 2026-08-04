package login

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var validate = validator.New()

// CreateLoginRequest is a model-mediator for transferring the gRPC request data to the Login service.
type CreateLoginRequest struct {
	Login       string
	Password    string
	UserID      uuid.UUID
	Title       string
	URL         string
	Description string
}

func (m *CreateLoginRequest) Validate() error {
	if m.Title == "" {
		return fmt.Errorf("title field is required")
	}
	return nil
}

// CreateLogin is a model-mediator for transferring the Login creation data to the Login repository.
type CreateLogin struct {
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;not null" validate:"required"`
	Login       string    `gorm:"column:hashed_login;not null" validate:"required"`
	Password    string    `gorm:"column:hashed_password;not null" validate:"required"`
	Title       string    `gorm:"column:title;not null" validate:"required"`
	Number      uint64    `gorm:"column:number;not null" validate:"required"`
	URL         string    `gorm:"column:url"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null" validate:"required"`
}

type Login struct {
	ID          uuid.UUID `gorm:"column:id"`
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid"`
	Login       string    `gorm:"column:hashed_login"`
	Password    string    `gorm:"column:hashed_password"`
	Title       string    `gorm:"column:title"`
	Number      uint64    `gorm:"column:number"`
	URL         string    `gorm:"column:url"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime"`
}

type UpdateLoginRequest struct {
	Title       *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=title,proto3" json:"title"`
	URL         *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
	Description *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Login       *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=login,proto3" json:"login" validate:"required"`
	Password    *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=password,proto3" json:"password"`
}

func (u *UpdateLoginRequest) Validate() error {
	if u.Title != nil && u.Title.String() == "" {
		return fmt.Errorf("title field is required")
	}
	return nil
}

type UpdateLogin struct {
	Login       string `gorm:"column:hashed_login"`
	Password    string `gorm:"column:hashed_password"`
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
