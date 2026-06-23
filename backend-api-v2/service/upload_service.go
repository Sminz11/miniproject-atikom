package service

import (
	"backend-api-v2/model"
	"backend-api-v2/repository"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type UploadService struct {
	Repo repository.UploadRepositoryInterface
}

func NewUploadService(repo repository.UploadRepositoryInterface) *UploadService {
	return &UploadService{Repo: repo}
}

func (s *UploadService) ProcessUpload(fileName string, fileBytes []byte) (*UploadResult, error) {
	scanner := bufio.NewScanner(bytes.NewReader(fileBytes))
	var details []struct {
		txnID    string
		amount   float64
		mobileNo string
	}
	hasHeader := false
	txnIDs := map[string]bool{}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")

		switch parts[0] {
		case "H":
			if len(parts) != 3 {
				return nil, errors.New("รูปแบบ Header ไม่ถูกต้อง")
			}
			hasHeader = true
		case "D":
			if len(parts) != 4 {
				return nil, errors.New("รูปแบบ Detail ไม่ถูกต้อง")
			}
			txnID := parts[1]
			amountStr := parts[2]
			mobileNo := parts[3]

			if txnIDs[txnID] {
				return nil, fmt.Errorf("Transaction ID ซ้ำ: %s", txnID)
			}
			txnIDs[txnID] = true

			amount, err := strconv.ParseFloat(amountStr, 64)
			if err != nil || amount <= 0 {
				return nil, fmt.Errorf("Amount ไม่ถูกต้อง: %s", amountStr)
			}

			if len(mobileNo) != 10 {
				return nil, fmt.Errorf("Mobile No ต้องมี 10 หลัก: %s", mobileNo)
			}

			details = append(details, struct {
				txnID    string
				amount   float64
				mobileNo string
			}{txnID, amount, mobileNo})
		}
	}

	if !hasHeader {
		return nil, errors.New("ไม่พบ Header ในไฟล์")
	}
	if len(details) == 0 {
		return nil, errors.New("ไม่พบข้อมูล Detail ในไฟล์")
	}

	uploadID, err := s.Repo.InsertHeader(fileName, len(details))
	if err != nil {
		return nil, fmt.Errorf("บันทึก Header ไม่สำเร็จ: %v", err)
	}

	for _, d := range details {
		if err := s.Repo.InsertDetail(uploadID, d.txnID, d.amount, d.mobileNo); err != nil {
			return nil, fmt.Errorf("บันทึก Detail ไม่สำเร็จ: %v", err)
		}
	}

	return &UploadResult{
		UploadID:    uploadID,
		FileName:    fileName,
		TotalRecord: len(details),
		Status:      "UPLOADED",
	}, nil
}

func (s *UploadService) GetHistory(fileName, status string, page, pageSize int) ([]model.UploadHeader, int, error) {
	return s.Repo.GetHeaders(fileName, status, page, pageSize)
}

func (s *UploadService) GetDetail(uploadID, page, pageSize int) ([]model.UploadDetail, int, error) {
	return s.Repo.GetDetails(uploadID, page, pageSize)
}
