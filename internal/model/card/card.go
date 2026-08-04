package card

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CreateCardRequest struct {
	PAN         *wrapperspb.UInt64Value `protobuf:"fixed64,1,opt,name=pan" json:"pan,omitempty"`
	Holder      *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=holder" json:"holder,omitempty"`
	ExpiryDate  *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=expiry_date" json:"expiry_date,omitempty"`
	CVV         *wrapperspb.UInt64Value `protobuf:"bytes,4,opt,name=cvv" json:"cvv,omitempty"`
	PIN         *wrapperspb.UInt64Value `protobuf:"bytes,5,opt,name=pin" json:"pin,omitempty"`
	Bank        *wrapperspb.StringValue `protobuf:"bytes,6,opt,name=bank" json:"bank,omitempty"`
	Brand       *wrapperspb.StringValue `protobuf:"bytes,7,opt,name=brand" json:"brand,omitempty"`
	UserID      uuid.UUID
	Title       *wrapperspb.StringValue `protobuf:"bytes,8,opt,name=title" json:"title,omitempty"`
	Description *wrapperspb.StringValue `protobuf:"bytes,9,opt,name=description" json:"description,omitempty"`
}

func (m *CreateCardRequest) Validate() error {
	if m.Title != nil && m.Title.String() == "" {
		return fmt.Errorf("title field is required")
	}
	return nil
}

type CreateCard struct {
	PAN         uint64 `gorm:"column:hashed_pan"`
	Holder      string `gorm:"column:hashed_holder"`
	ExpiryDate  string `gorm:"column:hashed_expiry_date"`
	CVV         uint64 `gorm:"column:hashed_cvv"`
	PIN         uint64 `gorm:"column:hashed_pin"`
	Bank        string
	Brand       string
	UserID      uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time
	Number      uint64
}

type Card struct {
	PAN         uint64 `gorm:"column:hashed_pan"`
	Holder      string `gorm:"column:hashed_holder"`
	ExpiryDate  string `gorm:"column:hashed_expiry_date"`
	CVV         uint64 `gorm:"column:hashed_cvv"`
	PIN         uint64 `gorm:"column:hashed_pin"`
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
	PAN         *wrapperspb.UInt64Value `protobuf:"fixed64,1,opt,name=pan,proto3" json:"pan"`
	Holder      *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=holder,proto3" json:"holder,omitempty"`
	ExpiryDate  *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=expiry_date,json=expiryDate,proto3" json:"expiry_date,omitempty"`
	CVV         *wrapperspb.UInt64Value `protobuf:"bytes,4,opt,name=cvv,proto3" json:"cvv,omitempty"`
	PIN         *wrapperspb.UInt64Value `protobuf:"bytes,5,opt,name=pin,proto3" json:"pin,omitempty"`
	Bank        *wrapperspb.StringValue `protobuf:"bytes,6,opt,name=bank,proto3" json:"bank,omitempty"`
	Brand       *wrapperspb.StringValue `protobuf:"bytes,7,opt,name=brand,proto3" json:"brand,omitempty"`
	Title       *wrapperspb.StringValue `protobuf:"bytes,8,opt,name=title,proto3" json:"title,omitempty"`
	Description *wrapperspb.StringValue `protobuf:"bytes,9,opt,name=description,proto3" json:"description,omitempty"`
}

func (m *UpdateCardRequest) Validate() error {
	if m.Title != nil && m.Title.String() == "" {
		return fmt.Errorf("title field is required")
	}
	return nil
}

type UpdateCard struct {
	PAN         uint64 `gorm:"column:hashed_pan"`
	Holder      string `gorm:"column:hashed_holder"`
	ExpiryDate  string `gorm:"column:hashed_expiry_date"`
	CVV         uint64 `gorm:"column:hashed_cvv"`
	PIN         uint64 `gorm:"column:hashed_pin"`
	Bank        string
	Brand       string
	Title       string
	Description string
	UpdatedAt   time.Time
}
