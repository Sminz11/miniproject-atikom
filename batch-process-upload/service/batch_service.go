package service

import (
	"batch-process-upload/repository"
	"fmt"
	"log"
)

type BatchService struct {
	Repo *repository.BatchRepository
}

func NewBatchService(repo *repository.BatchRepository) *BatchService {
	return &BatchService{Repo: repo}
}

func (s *BatchService) ProcessAll() error {
	log.Println("Batch เริ่มทำงาน...")

	headers, err := s.Repo.GetPendingHeaders()
	if err != nil {
		return fmt.Errorf("ดึง Header ไม่ได้: %v", err)
	}

	if len(headers) == 0 {
		log.Println("ไม่มีรายการที่ต้องประมวลผล")
		return nil
	}

	for _, h := range headers {
		log.Printf("กำลังประมวลผล Upload ID: %d, File: %s", h.UploadID, h.FileName)

		// อัปเดตเป็น PROCESSING
		s.Repo.UpdateHeaderProcessing(h.UploadID)

		// ดึง Detail
		details, err := s.Repo.GetPendingDetails(h.UploadID)
		if err != nil {
			log.Printf("ดึง Detail ไม่ได้ Upload ID %d: %v", h.UploadID, err)
			continue
		}

		totalSuccess := 0
		totalFailed := 0

		for _, d := range details {
			status, errMsg := processMockLogic(d.Amount, d.MobileNo, d.TransactionID)

			err := s.Repo.UpdateDetail(d.DetailID, status, errMsg)
			if err != nil {
				log.Printf("อัปเดต Detail ID %d ไม่ได้: %v", d.DetailID, err)
				continue
			}

			if status == "SUCCESS" {
				totalSuccess++
			} else {
				totalFailed++
			}

			log.Printf("  TxnID: %s → %s %s", d.TransactionID, status, errMsg)
		}

		// อัปเดต Header เป็น COMPLETED
		s.Repo.UpdateHeaderCompleted(h.UploadID, totalSuccess, totalFailed)
		log.Printf("Upload ID %d เสร็จ: SUCCESS=%d, FAILED=%d", h.UploadID, totalSuccess, totalFailed)
	}

	log.Println("Batch ทำงานเสร็จสิ้น")
	return nil
}

// Mock Business Logic ตาม Requirement
func processMockLogic(amount float64, mobileNo, txnID string) (string, string) {
	if txnID == "" {
		return "FAILED", "Transaction ID is required"
	}
	if len(mobileNo) != 10 {
		return "FAILED", "Invalid mobile number"
	}
	if amount > 300 {
		return "FAILED", "Amount exceed limit"
	}
	return "SUCCESS", ""
}
