package repository

import (
	"batch-process-upload/model"
	"database/sql"
	"fmt"
	"time"
)

type BatchRepository struct {
	DB     *sql.DB
	Schema string
}

func NewBatchRepository(db *sql.DB, schema string) *BatchRepository {
	return &BatchRepository{DB: db, Schema: schema}
}

func (r *BatchRepository) GetPendingHeaders() ([]model.UploadHeader, error) {
	query := fmt.Sprintf(`SELECT upload_id, file_name, status, total_record 
		FROM %s.upload_header WHERE status = 'UPLOADED'`, r.Schema)
	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var headers []model.UploadHeader
	for rows.Next() {
		var h model.UploadHeader
		rows.Scan(&h.UploadID, &h.FileName, &h.Status, &h.TotalRecord)
		headers = append(headers, h)
	}
	return headers, nil
}

func (r *BatchRepository) UpdateHeaderProcessing(uploadID int) error {
	query := fmt.Sprintf(`UPDATE %s.upload_header SET status = 'PROCESSING', updated_date = $1 
		WHERE upload_id = $2`, r.Schema)
	_, err := r.DB.Exec(query, time.Now(), uploadID)
	return err
}

func (r *BatchRepository) GetPendingDetails(uploadID int) ([]model.UploadDetail, error) {
	query := fmt.Sprintf(`SELECT detail_id, upload_id, transaction_id, amount, mobile_no, status 
		FROM %s.upload_detail WHERE upload_id = $1 AND status = 'PENDING'`, r.Schema)
	rows, err := r.DB.Query(query, uploadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details []model.UploadDetail
	for rows.Next() {
		var d model.UploadDetail
		rows.Scan(&d.DetailID, &d.UploadID, &d.TransactionID, &d.Amount, &d.MobileNo, &d.Status)
		details = append(details, d)
	}
	return details, nil
}

func (r *BatchRepository) UpdateDetail(detailID int, status, errMsg string) error {
	query := fmt.Sprintf(`UPDATE %s.upload_detail SET status = $1, error_message = $2, 
		process_date = $3, updated_date = $3 WHERE detail_id = $4`, r.Schema)
	_, err := r.DB.Exec(query, status, errMsg, time.Now(), detailID)
	return err
}

func (r *BatchRepository) UpdateHeaderCompleted(uploadID, totalSuccess, totalFailed int) error {
	query := fmt.Sprintf(`UPDATE %s.upload_header SET status = 'COMPLETED', 
		total_success = $1, total_failed = $2, updated_date = $3 WHERE upload_id = $4`, r.Schema)
	_, err := r.DB.Exec(query, totalSuccess, totalFailed, time.Now(), uploadID)
	return err
}
