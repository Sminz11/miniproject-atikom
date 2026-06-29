package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type RetryRepository struct {
	DB     *sql.DB
	Schema string
}

func NewRetryRepository(db *sql.DB, schema string) *RetryRepository {
	return &RetryRepository{DB: db, Schema: schema}
}

func (r *RetryRepository) RetryDetail(uploadID, detailID int) error {
	// ตรวจสอบว่า Header ไม่ได้อยู่ใน PROCESSING
	var headerStatus string
	err := r.DB.QueryRow(fmt.Sprintf(
		`SELECT status FROM %s.upload_header WHERE upload_id = $1`, r.Schema),
		uploadID).Scan(&headerStatus)
	if err != nil {
		return fmt.Errorf("ไม่พบ Upload ID: %d", uploadID)
	}
	if headerStatus == "PROCESSING" {
		return fmt.Errorf("ไม่สามารถ Retry ขณะ Header อยู่ในสถานะ PROCESSING")
	}

	// ตรวจสอบว่า Detail เป็น FAILED
	var detailStatus string
	err = r.DB.QueryRow(fmt.Sprintf(
		`SELECT status FROM %s.upload_detail WHERE detail_id = $1 AND upload_id = $2`,
		r.Schema), detailID, uploadID).Scan(&detailStatus)
	if err != nil {
		return fmt.Errorf("ไม่พบ Detail ID: %d", detailID)
	}
	if detailStatus != "FAILED" {
		return fmt.Errorf("Retry ได้เฉพาะรายการที่ status = FAILED เท่านั้น")
	}

	// เปลี่ยน status เป็น PENDING และล้าง error_message
	_, err = r.DB.Exec(fmt.Sprintf(
		`UPDATE %s.upload_detail SET status = 'PENDING', error_message = '', 
		process_date = NULL, updated_date = $1 
		WHERE detail_id = $2`, r.Schema), time.Now(), detailID)
	if err != nil {
		return fmt.Errorf("Retry ไม่สำเร็จ: %v", err)
	}

	// เปลี่ยน Header กลับเป็น UPLOADED เพื่อให้ Batch ประมวลผลใหม่
	_, err = r.DB.Exec(fmt.Sprintf(
		`UPDATE %s.upload_header SET status = 'UPLOADED', updated_date = $1 
		WHERE upload_id = $2`, r.Schema), time.Now(), uploadID)
	return err
}

func (r *RetryRepository) GetDetailsForExport(uploadID int) ([]map[string]interface{}, error) {
	rows, err := r.DB.Query(fmt.Sprintf(
		`SELECT detail_id, transaction_id, amount, mobile_no, status, 
		COALESCE(error_message,''), process_date, created_date
		FROM %s.upload_detail WHERE upload_id = $1 ORDER BY detail_id`, r.Schema),
		uploadID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []map[string]interface{}
	for rows.Next() {
		var detailID int
		var txnID, mobileNo, status, errMsg string
		var amount float64
		var processDate *time.Time
		var createdDate time.Time

		rows.Scan(&detailID, &txnID, &amount, &mobileNo, &status,
			&errMsg, &processDate, &createdDate)

		results = append(results, map[string]interface{}{
			"detailId":      detailID,
			"transactionId": txnID,
			"amount":        amount,
			"mobileNo":      mobileNo,
			"status":        status,
			"errorMessage":  errMsg,
			"processDate":   processDate,
			"createdDate":   createdDate,
		})
	}
	return results, nil
}
