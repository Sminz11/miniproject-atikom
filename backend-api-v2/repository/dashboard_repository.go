package repository

import (
	"backend-api-v2/model"
	"database/sql"
	"fmt"
)

type DashboardRepository struct {
	DB     *sql.DB
	Schema string
}

func NewDashboardRepository(db *sql.DB, schema string) *DashboardRepository {
	return &DashboardRepository{DB: db, Schema: schema}
}

func (r *DashboardRepository) GetSummary() (*model.DashboardSummary, error) {
	summary := &model.DashboardSummary{
		UploadByStatus: make(map[string]int),
	}

	// Upload by status
	rows, err := r.DB.Query(fmt.Sprintf(
		`SELECT status, COUNT(*) FROM %s.upload_header GROUP BY status`, r.Schema))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		rows.Scan(&status, &count)
		summary.UploadByStatus[status] = count
		summary.TotalUpload += count
	}

	// Detail summary
	r.DB.QueryRow(fmt.Sprintf(
		`SELECT COUNT(*) FROM %s.upload_detail WHERE status = 'SUCCESS'`, r.Schema)).
		Scan(&summary.TotalSuccess)
	r.DB.QueryRow(fmt.Sprintf(
		`SELECT COUNT(*) FROM %s.upload_detail WHERE status = 'FAILED'`, r.Schema)).
		Scan(&summary.TotalFailed)
	r.DB.QueryRow(fmt.Sprintf(
		`SELECT COUNT(*) FROM %s.upload_detail WHERE status = 'PENDING'`, r.Schema)).
		Scan(&summary.TotalPending)

	// Recent uploads
	recentRows, err := r.DB.Query(fmt.Sprintf(
		`SELECT upload_id, file_name, upload_date, total_record, total_success, total_failed, status, created_date
		FROM %s.upload_header ORDER BY upload_id DESC LIMIT 5`, r.Schema))
	if err != nil {
		return nil, err
	}
	defer recentRows.Close()
	for recentRows.Next() {
		var h model.UploadHeader
		recentRows.Scan(&h.UploadID, &h.FileName, &h.UploadDate, &h.TotalRecord,
			&h.TotalSuccess, &h.TotalFailed, &h.Status, &h.CreatedDate)
		summary.RecentUploads = append(summary.RecentUploads, h)
	}

	return summary, nil
}
