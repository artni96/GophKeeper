package text

import "time"

type Text struct {
	Text        string
	Title       string
	Number      uint64
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
