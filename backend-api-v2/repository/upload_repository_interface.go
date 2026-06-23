package repository

import "backend-api-v2/model"

type UploadRepositoryInterface interface {
	InsertHeader(fileName string, totalRecord int) (int, error)
	InsertDetail(uploadID int, txnID string, amount float64, mobileNo string) error
	GetHeaders(fileName, status string, page, pageSize int) ([]model.UploadHeader, int, error)
	GetDetails(uploadID, page, pageSize int) ([]model.UploadDetail, int, error)
}
