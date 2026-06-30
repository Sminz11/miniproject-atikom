package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type AuditLogRepository struct {
	DB     *sql.DB
	Schema string
}

func NewAuditLogRepository(db *sql.DB, schema string) *AuditLogRepository {
	return &AuditLogRepository{DB: db, Schema: schema}
}

func (r *AuditLogRepository) Insert(username, action, refID, description string) error {
	query := fmt.Sprintf(`INSERT INTO %s.audit_log 
		(username, action, ref_id, description, created_date) 
		VALUES ($1, $2, $3, $4, $5)`, r.Schema)
	_, err := r.DB.Exec(query, username, action, refID, description, time.Now())
	return err
}
