package repository

import (
	"backend-api-v2/model"
	"database/sql"
	"fmt"
	"time"
)

type UploadRepository struct {
	DB     *sql.DB
	Schema string
}

func NewUploadRepository(db *sql.DB, schema string) *UploadRepository {
	return &UploadRepository{DB: db, Schema: schema}
}

func (r *UploadRepository) InsertHeader(fileName string, totalRecord int) (int, error) {
	var uploadID int
	query := fmt.Sprintf(`INSERT INTO %s.upload_header 
		(file_name, upload_date, total_record, status, created_date) 
		VALUES ($1, $2, $3, 'UPLOADED', $2) RETURNING upload_id`, r.Schema)
	err := r.DB.QueryRow(query, fileName, time.Now(), totalRecord).Scan(&uploadID)
	return uploadID, err
}

func (r *UploadRepository) InsertDetail(uploadID int, txnID string, amount float64, mobileNo string) error {
	query := fmt.Sprintf(`INSERT INTO %s.upload_detail 
		(upload_id, transaction_id, amount, mobile_no, status, created_date) 
		VALUES ($1, $2, $3, $4, 'PENDING', $5)`, r.Schema)
	_, err := r.DB.Exec(query, uploadID, txnID, amount, mobileNo, time.Now())
	return err
}

func (r *UploadRepository) GetHeaders(fileName, status string, page, pageSize int) ([]model.UploadHeader, int, error) {
	offset := (page - 1) * pageSize
	query := fmt.Sprintf(`SELECT upload_id, file_name, upload_date, total_record, 
		total_success, total_failed, status, created_date
		FROM %s.upload_header WHERE 1=1`, r.Schema)
	args := []interface{}{}
	i := 1

	if fileName != "" {
		query += fmt.Sprintf(" AND file_name ILIKE $%d", i)
		args = append(args, "%"+fileName+"%")
		i++
	}
	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", i)
		args = append(args, status)
		i++
	}

	var total int
	countQuery := "SELECT COUNT(*) FROM (" + query + ") t"
	r.DB.QueryRow(countQuery, args...).Scan(&total)

	query += fmt.Sprintf(" ORDER BY upload_id DESC LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, pageSize, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var headers []model.UploadHeader
	for rows.Next() {
		var h model.UploadHeader
		rows.Scan(&h.UploadID, &h.FileName, &h.UploadDate, &h.TotalRecord,
			&h.TotalSuccess, &h.TotalFailed, &h.Status, &h.CreatedDate)
		headers = append(headers, h)
	}
	return headers, total, nil
}

func (r *UploadRepository) GetDetails(uploadID, page, pageSize int) ([]model.UploadDetail, int, error) {
	offset := (page - 1) * pageSize
	var total int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s.upload_detail WHERE upload_id = $1", r.Schema)
	r.DB.QueryRow(countQuery, uploadID).Scan(&total)

	query := fmt.Sprintf(`SELECT detail_id, upload_id, transaction_id, amount, mobile_no, 
		status, COALESCE(error_message,''), process_date, created_date
		FROM %s.upload_detail WHERE upload_id = $1 ORDER BY detail_id LIMIT $2 OFFSET $3`, r.Schema)

	rows, err := r.DB.Query(query, uploadID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var details []model.UploadDetail
	for rows.Next() {
		var d model.UploadDetail
		rows.Scan(&d.DetailID, &d.UploadID, &d.TransactionID, &d.Amount,
			&d.MobileNo, &d.Status, &d.ErrorMessage, &d.ProcessDate, &d.CreatedDate)
		details = append(details, d)
	}
	return details, total, nil
}
