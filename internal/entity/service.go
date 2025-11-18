package entity

import "time"

type Service struct {
	ID                int       `json:"id"`
	Name              string    `json:"name" binding:"required"`
	Description       *string   `json:"description"`
	SpecificPhoto     *string   `json:"specific_photo"`
	Price             *int      `json:"price"`
	ServiceCategoryID *int      `json:"service_category_id"`
	SpecializationID  *int      `json:"specialization_id"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type ServiceCreateRequest struct {
	Name              string  `json:"name" binding:"required"`
	Description       *string `json:"description"`
	SpecificPhoto     *string `json:"specific_photo"`
	Price             *int    `json:"price"`
	ServiceCategoryID *int    `json:"service_category_id"`
	SpecializationID  *int    `json:"specialization_id"`
}
