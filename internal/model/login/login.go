package login

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
	var errs []string
	if m.Login == "" {
		errs = append(errs, "login field is required")
	}
	if m.Password == "" {
		errs = append(errs, "password field is required")
	}
	if m.Title == "" {
		errs = append(errs, "title field is required")
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, ", "))
	}
	return nil
}

// CreateLogin is a model-mediator for transferring the Login creation data to the Login repository.
type CreateLogin struct {
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid;not null" validate:"required"`
	Login       string    `gorm:"column:hashed_login;not null" validate:"required"`
	Password    string    `gorm:"column:hashed_password;not null" validate:"required"`
	Title       string    `gorm:"column:title;not null" validate:"required"`
	Number      int64     `gorm:"column:number;not null" validate:"required"`
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
	Number      int64     `gorm:"column:number"`
	URL         string    `gorm:"column:url"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime"`
}

type UpdateLoginRequest struct {
	Login       string `gorm:"column:hashed_login"`
	Password    string `gorm:"column:hashed_password"`
	Title       string `gorm:"column:title"`
	URL         string `gorm:"column:url"`
	Description string `gorm:"column:description"`
}

type UpdateLogin struct {
	Login       string    `gorm:"column:hashed_login"`
	Password    string    `gorm:"column:hashed_password"`
	Title       string    `gorm:"column:title"`
	URL         string    `gorm:"column:url"`
	Description string    `gorm:"column:description"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:datetime"`
}

type GetListLoginResponse struct {
	Title       string
	Description string
	Number      int64
}
