package model

import "time"

type UploadHeader struct {
	UploadID     int        `json:"uploadId"`
	FileName     string     `json:"fileName"`
	Status       string     `json:"status"`
	TotalRecord  int        `json:"totalRecord"`
	TotalSuccess int        `json:"totalSuccess"`
	TotalFailed  int        `json:"totalFailed"`
	CreatedDate  time.Time  `json:"createdDate"`
	UpdatedDate  *time.Time `json:"updatedDate"`
}

type UploadDetail struct {
	DetailID      int        `json:"detailId"`
	UploadID      int        `json:"uploadId"`
	TransactionID string     `json:"transactionId"`
	Amount        float64    `json:"amount"`
	MobileNo      string     `json:"mobileNo"`
	Status        string     `json:"status"`
	ErrorMessage  string     `json:"errorMessage"`
	ProcessDate   *time.Time `json:"processDate"`
}
