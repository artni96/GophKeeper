package card

import (
	"time"
)

type Card struct {
	PAN         uint64
	Holder      string
	ExpiryDate  string
	CVV         string
	PIN         string
	Bank        string
	Brand       string
	Title       string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Number      uint64
}
