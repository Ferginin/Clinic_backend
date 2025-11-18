package entity

import "time"

type Doctor struct {
	ID              int              `json:"id"`
	Fullname        string           `json:"fullname" binding:"required"`
	Description     *string          `json:"description"`
	DoctorPhoto     *string          `json:"doctor_photo"`
	ScheduleID      *int             `json:"schedule_id"`
	Schedule        *Schedule        `json:"schedule,omitempty"`
	Specializations []Specialization `json:"specializations,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}

type DoctorCreateRequest struct {
	Fullname           string  `json:"fullname" binding:"required"`
	Description        *string `json:"description"`
	DoctorPhoto        *string `json:"doctor_photo"`
	ScheduleID         *int    `json:"schedule_id"`
	SpecializationIDs  []int   `json:"specialization_ids"`
}

type DoctorUpdateRequest struct {
	Fullname           *string `json:"fullname"`
	Description        *string `json:"description"`
	DoctorPhoto        *string `json:"doctor_photo"`
	ScheduleID         *int    `json:"schedule_id"`
	SpecializationIDs  []int   `json:"specialization_ids"`
}
