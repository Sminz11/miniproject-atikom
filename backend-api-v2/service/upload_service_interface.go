package service

import "backend-api-v2/model"

type UploadServiceInterface interface {
	ProcessUpload(fileName string, fileBytes []byte) (*UploadResult, error)
	GetHistory(fileName, status string, page, pageSize int) ([]model.UploadHeader, int, error)
	GetDetail(uploadID, page, pageSize int) ([]model.UploadDetail, int, error)
}

type UploadResult struct {
	UploadID    int    `json:"uploadId"`
	FileName    string `json:"fileName"`
	TotalRecord int    `json:"totalRecord"`
	Status      string `json:"status"`
}
