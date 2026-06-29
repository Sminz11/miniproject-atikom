package repository

import (
	"backend-api-v2/model"
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

func (r *AuditLogRepository) GetLogs(action string, page, pageSize int) ([]model.AuditLog, int, error) {
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`SELECT audit_id, COALESCE(username,''), action, 
		COALESCE(ref_id,''), COALESCE(description,''), created_date
		FROM %s.audit_log WHERE 1=1`, r.Schema)
	args := []interface{}{}
	i := 1

	if action != "" {
		query += fmt.Sprintf(" AND action = $%d", i)
		args = append(args, action)
		i++
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM (" + query + ") t"
	r.DB.QueryRow(countQuery, args...).Scan(&total)

	query += fmt.Sprintf(" ORDER BY audit_id DESC LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, pageSize, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var logs []model.AuditLog
	for rows.Next() {
		var l model.AuditLog
		rows.Scan(&l.AuditID, &l.Username, &l.Action, &l.RefID, &l.Description, &l.CreatedDate)
		logs = append(logs, l)
	}
	return logs, total, nil
}
