package login

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// CreateLoginRequest is a model-mediator for transferring the gRPC request data to the Login service.
type CreateLoginRequest struct {
	Login       *wrapperspb.StringValue
	Password    *wrapperspb.StringValue
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
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;not null"`
	Login       string    `gorm:"column:hashed_login;not null"`
	Password    string    `gorm:"column:hashed_password;not null"`
	Title       string    `gorm:"column:title;not null"`
	Number      uint64    `gorm:"column:number;not null"`
	URL         string    `gorm:"column:url"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null"`
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
	Login       *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=login,proto3" json:"login,omitempty"`
	Password    *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
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
