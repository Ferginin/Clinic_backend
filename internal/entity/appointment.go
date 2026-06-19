package entity

import "time"

const AppointmentDurationMinutes = 30

type Appointment struct {
	ID          int       `json:"id"`
	PatientID   int       `json:"patient_id" binding:"required"`
	DoctorID    int       `json:"doctor_id" binding:"required"`
	Doctor      *Doctor   `json:"doctor,omitempty"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
	Status      string    `json:"status"`
	Notes       *string   `json:"notes,omitempty"`
	Results     *string   `json:"results,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AppointmentCreateRequest struct {
	DoctorID    int     `json:"doctor_id" binding:"required"`
	ScheduledAt string  `json:"scheduled_at" binding:"required"`
	Notes       *string `json:"notes,omitempty"`
}

type AppointmentUpdateRequest struct {
	ScheduledAt *string `json:"scheduled_at,omitempty"`
	Status      *string `json:"status,omitempty"`
	Notes       *string `json:"notes,omitempty"`
}

type AppointmentResultRequest struct {
	Results string `json:"results" binding:"required"`
	Status  string `json:"status" binding:"required,oneof=pending confirmed completed canceled"`
}

// AppointmentAdminCreateRequest used by admin to book on behalf of a patient or guest
type AppointmentAdminCreateRequest struct {
	DoctorID    int     `json:"doctor_id" binding:"required"`
	ScheduledAt string  `json:"scheduled_at" binding:"required"`
	Notes       *string `json:"notes,omitempty"`
	// Existing patient mode
	PatientID *int `json:"patient_id,omitempty"`
	// Guest (callback) mode
	GuestName  *string `json:"guest_name,omitempty"`
	GuestPhone *string `json:"guest_phone,omitempty"`
}

// TimeSlot представляет доступный слот для записи
type TimeSlot struct {
	StartTime string `json:"start_time"` // "09:00"
	EndTime   string `json:"end_time"`   // "09:30"
	Available bool   `json:"available"`
}

// AvailableSlotsResponse ответ с доступными слотами
type AvailableSlotsResponse struct {
	DoctorID int        `json:"doctor_id"`
	Date     string     `json:"date"`
	Slots    []TimeSlot `json:"slots"`
}

// ValidStatusTransitions допустимые переходы статусов
var ValidStatusTransitions = map[string][]string{
	"pending":   {"confirmed", "canceled"},
	"confirmed": {"completed", "canceled"},
	"completed": {},
	"canceled":  {},
}

// IsValidTransition проверяет допустимость перехода статуса
func IsValidTransition(from, to string) bool {
	allowed, exists := ValidStatusTransitions[from]
	if !exists {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
