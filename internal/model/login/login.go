package login

import (
	"time"

	"github.com/google/uuid"
)

// CreateLoginRequest is a model-mediator for transferring the gRPC request data to the Login service.
type CreateLoginRequest struct {
	Login       string
	Password    string
	UserID      uuid.UUID
	Title       string
	URL         string
	Description string
}

// CreateLogin is a model-mediator for transferring the Login creation data to the Login repository.
type CreateLogin struct {
	UserID      uuid.UUID `gorm:"column:user_id;type:uuid"`
	Login       string    `gorm:"column:hashed_login"`
	Password    string    `gorm:"column:hashed_password"`
	Title       string    `gorm:"column:title"`
	Number      int64     `gorm:"column:number"`
	URL         string    `gorm:"column:url"`
	Description string    `gorm:"column:description"`
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime"`
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
