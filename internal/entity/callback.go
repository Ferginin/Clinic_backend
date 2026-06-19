package entity

import "time"

type CallbackRequest struct {
	ID        int       `json:"id"`
	Name      string    `json:"name" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	Message   *string   `json:"message,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CallbackRequestCreate struct {
	Name    string  `json:"name" binding:"required,min=2"`
	Phone   string  `json:"phone" binding:"required"`
	Message *string `json:"message,omitempty"`
}

type CallbackRequestUpdate struct {
	Status *string `json:"status,omitempty"`
}
