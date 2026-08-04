package card

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var ErrInvalidPAN = errors.New("invalid PAN value")
var ErrInvalidCVV = errors.New("invalid CVV value")
var ErrInvalidPIN = errors.New("invalid PIN value")

type CreateCardRequest struct {
	PAN         *wrapperspb.UInt64Value `protobuf:"fixed64,1,opt,name=pan" json:"pan,omitempty"`
	Holder      *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=holder" json:"holder,omitempty"`
	ExpiryDate  *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=expiry_date" json:"expiry_date,omitempty"`
	CVV         *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=cvv" json:"cvv,omitempty"`
	PIN         *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=pin" json:"pin,omitempty"`
	Bank        *wrapperspb.StringValue `protobuf:"bytes,6,opt,name=bank" json:"bank,omitempty"`
	Brand       *wrapperspb.StringValue `protobuf:"bytes,7,opt,name=brand" json:"brand,omitempty"`
	UserID      uuid.UUID
	Title       *wrapperspb.StringValue `protobuf:"bytes,8,opt,name=title" json:"title,omitempty"`
	Description *wrapperspb.StringValue `protobuf:"bytes,9,opt,name=description" json:"description,omitempty"`
}

func (c *CreateCardRequest) Validate() error {
	err := cardValidator(c.PAN, c.CVV, c.PIN, c.Title)
	if err != nil {
		return err
	}
	return nil
}

type CreateCard struct {
	PAN         uint64 `gorm:"column:hashed_pan"`
	Holder      string `gorm:"column:hashed_holder"`
	ExpiryDate  string `gorm:"column:hashed_expiry_date"`
	CVV         string `gorm:"column:hashed_cvv"`
	PIN         string `gorm:"column:hashed_pin"`
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
	CVV         string `gorm:"column:hashed_cvv"`
	PIN         string `gorm:"column:hashed_pin"`
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
	CVV         *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=cvv,proto3" json:"cvv,omitempty"`
	PIN         *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=pin,proto3" json:"pin,omitempty"`
	Bank        *wrapperspb.StringValue `protobuf:"bytes,6,opt,name=bank,proto3" json:"bank,omitempty"`
	Brand       *wrapperspb.StringValue `protobuf:"bytes,7,opt,name=brand,proto3" json:"brand,omitempty"`
	Title       *wrapperspb.StringValue `protobuf:"bytes,8,opt,name=title,proto3" json:"title,omitempty"`
	Description *wrapperspb.StringValue `protobuf:"bytes,9,opt,name=description,proto3" json:"description,omitempty"`
}

func (c *UpdateCardRequest) Validate() error {
	err := cardValidator(c.PAN, c.CVV, c.PIN, c.Title)
	if err != nil {
		return err
	}
	return nil
}

type UpdateCard struct {
	PAN         uint64 `gorm:"column:hashed_pan"`
	Holder      string `gorm:"column:hashed_holder"`
	ExpiryDate  string `gorm:"column:hashed_expiry_date"`
	CVV         string `gorm:"column:hashed_cvv"`
	PIN         string `gorm:"column:hashed_pin"`
	Bank        string
	Brand       string
	Title       string
	Description string
	UpdatedAt   time.Time
}

// luhnValidator checks if a card number is valid or not.
func luhnValidator(number uint64) bool {
	strNumber := strconv.FormatUint(number, 10)
	sum := 0
	strLength := len(strNumber)
	parity := strLength % 2

	for i := 0; i < strLength-1; i++ {
		digit := int(strNumber[i] - '0')
		if parity == (i+1)%2 {
			sum += digit
		} else if digit > 4 {
			sum += 2*digit - 9
		} else {
			sum += 2 * digit
		}
	}
	controlSum := (10 - (sum % 10)) % 10
	return int(strNumber[strLength-1]-'0') == controlSum
}

// cardValidator unifies the single validation flow for UpdateCardRequest and CreateCardRequest.
func cardValidator(pan *wrapperspb.UInt64Value, cvv, pin, title *wrapperspb.StringValue) error {
	if title != nil && title.String() == "" {
		return fmt.Errorf("title field is required")
	}
	if pan != nil && pan.Value != 0 {
		ok := luhnValidator(pan.Value)
		if !ok {
			return ErrInvalidPAN
		}
	}
	if cvv != nil {
		strCVV := cvv.Value
		intCVV, err := strconv.Atoi(strCVV)
		if err != nil {
			return err
		}
		if (len(strCVV) == 4 || len(strCVV) == 3) && intCVV > 0 {
			return nil
		}
		return fmt.Errorf("invalid CVV value: %d", intCVV)
	}

	if pin != nil {
		strPIN := pin.Value
		intPIN, err := strconv.Atoi(strPIN)
		if err != nil {
			return err
		}
		if len(strPIN) != 4 && (intPIN < 0 || intPIN > 10000) {
			return fmt.Errorf("invalid PIN value: %d", intPIN)
		}
	}
	return nil
}
