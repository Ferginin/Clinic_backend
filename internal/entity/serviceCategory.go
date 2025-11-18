package entity

import "time"

type ServiceCategory struct {
	ID               int       `json:"id"`
	Name             string    `json:"name" binding:"required"`
	Description      *string   `json:"description"`
	CategoryPhoto    *string   `json:"category_photo"`
	Favorite         bool      `json:"favorite"`
	SpecializationID *int      `json:"specialization_id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
