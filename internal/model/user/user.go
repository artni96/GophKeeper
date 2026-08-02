package user

import "github.com/google/uuid"

type UserCreateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserCreate struct {
	Username       string `gorm:"column:username"`
	HashedPassword string `gorm:"column:hashed_password"`
}

type LoginRequest struct {
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}

type User struct {
	ID             uuid.UUID `gorm:"column:id;type:uuid;"`
	Username       string    `gorm:"column:username"`
	HashedPassword string    `gorm:"column:hashed_password"`
}
