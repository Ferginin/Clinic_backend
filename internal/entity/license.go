package entity

import "time"

type License struct {
	ID          int       `json:"id"`
	Photo       *string   `json:"photo"`
	Name        *string   `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
