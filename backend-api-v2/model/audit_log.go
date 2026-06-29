package model

import "time"

type AuditLog struct {
	AuditID     int       `json:"auditId"`
	Username    string    `json:"username"`
	Action      string    `json:"action"`
	RefID       string    `json:"refId"`
	Description string    `json:"description"`
	CreatedDate time.Time `json:"createdDate"`
}
