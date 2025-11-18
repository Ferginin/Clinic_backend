package entity

import "time"

type Schedule struct {
	ID        int       `json:"id"`
	Day       int       `json:"day" binding:"required,min=1,max=7"` // 1=Monday, 7=Sunday
	TimeFrom  string    `json:"time_from" binding:"required"`
	TimeTo    string    `json:"time_to" binding:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
