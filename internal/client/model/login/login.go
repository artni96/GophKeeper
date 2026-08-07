package login

import (
	"time"
)

type Login struct {
	Login       string
	Password    string
	Title       string
	Number      uint64
	URL         string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
