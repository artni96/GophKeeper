package card

import (
	"time"

	"github.com/artni96/GophKeeper/internal/server/model/common"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CreateCardRequest struct {
	PAN        common.PBEncryptedField
	Holder     common.PBEncryptedField
	ExpiryDate common.PBEncryptedField
	CVV        common.PBEncryptedField
	PIN        common.PBEncryptedField

	Bank        *wrapperspb.StringValue
	Brand       *wrapperspb.StringValue
	UserID      uuid.UUID
	Title       *wrapperspb.StringValue
	Description *wrapperspb.StringValue
}

type CreateCard struct {
	PAN        common.EncryptedField `gorm:"embedded;embeddedPrefix:pan_"`
	Holder     common.EncryptedField `gorm:"embedded;embeddedPrefix:holder_"`
	ExpiryDate common.EncryptedField `gorm:"embedded;embeddedPrefix:expiry_date_"`
	CVV        common.EncryptedField `gorm:"embedded;embeddedPrefix:cvv_"`
	PIN        common.EncryptedField `gorm:"embedded;embeddedPrefix:pin_"`

	Bank        string
	Brand       string
	UserID      uuid.UUID
	Title       string
	Description string
	CreatedAt   time.Time
	Number      uint64
}

func (c *CreateCard) GetUserID() uuid.UUID {
	return c.UserID
}

func (c *CreateCard) SetNumber(number uint64) {
	c.Number = number
}

type Card struct {
	ID         uuid.UUID
	PAN        common.EncryptedField `gorm:"embedded;embeddedPrefix:pan_"`
	Holder     common.EncryptedField `gorm:"embedded;embeddedPrefix:holder_"`
	ExpiryDate common.EncryptedField `gorm:"embedded;embeddedPrefix:expiry_date_"`
	CVV        common.EncryptedField `gorm:"embedded;embeddedPrefix:cvv_"`
	PIN        common.EncryptedField `gorm:"embedded;embeddedPrefix:pin_"`

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
	PAN        common.PBEncryptedField
	Holder     common.PBEncryptedField
	ExpiryDate common.PBEncryptedField
	CVV        common.PBEncryptedField
	PIN        common.PBEncryptedField

	Bank        *wrapperspb.StringValue
	Brand       *wrapperspb.StringValue
	Title       *wrapperspb.StringValue
	Description *wrapperspb.StringValue
}

type UpdateCard struct {
	PAN        common.EncryptedField `gorm:"embedded;embeddedPrefix:pan_"`
	Holder     common.EncryptedField `gorm:"embedded;embeddedPrefix:holder_"`
	ExpiryDate common.EncryptedField `gorm:"embedded;embeddedPrefix:expiry_date_"`
	CVV        common.EncryptedField `gorm:"embedded;embeddedPrefix:cvv_"`
	PIN        common.EncryptedField `gorm:"embedded;embeddedPrefix:pin_"`

	Bank        string
	Brand       string
	Title       string
	Description string
	UpdatedAt   time.Time
}
