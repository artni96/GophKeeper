package text

import (
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CreateTextRequest struct {
	Title       *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=title" json:"title,omitempty"`
	Description *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	Text        *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=text" json:"text,omitempty"`
	UserID      uuid.UUID
}

type CreateText struct {
	Title       string
	Description string
	Text        string `gorm:"column:hashed_text"`
	UserID      uuid.UUID
	CreatedAt   time.Time
	Number      uint64
}

type Text struct {
	Title       string
	Description string
	Text        string `gorm:"column:hashed_text"`
	UserID      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Number      uint64
}

type UpdateTextRequest struct {
	Title       *wrapperspb.StringValue `protobuf:"bytes,1,opt,name=title" json:"title,omitempty"`
	Description *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=description" json:"description,omitempty"`
	Text        *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=text" json:"text,omitempty"`
}

type UpdateText struct {
	Title       string
	Description string
	Text        string `gorm:"column:hashed_text"`
	UpdatedAt   time.Time
}
