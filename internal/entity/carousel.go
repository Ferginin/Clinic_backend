package entity

import "time"

type Carousel struct {
	ID          int       `json:"id"`
	Image       *string   `json:"image"`
	Header      *string   `json:"header"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
