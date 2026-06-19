package entity

import "time"

type AuditLog struct {
	ID         int       `json:"id"`
	AdminID    int       `json:"admin_id"`
	Action     string    `json:"action"`
	EntityType string    `json:"entity_type"`
	EntityID   *int      `json:"entity_id,omitempty"`
	IP         string    `json:"ip"`
	CreatedAt  time.Time `json:"created_at"`
}
