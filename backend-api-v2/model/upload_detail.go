package model

import "time"

type UploadDetail struct {
	DetailID      int        `json:"detailId"`
	UploadID      int        `json:"uploadId"`
	TransactionID string     `json:"transactionId"`
	Amount        float64    `json:"amount"`
	MobileNo      string     `json:"mobileNo"`
	Status        string     `json:"status"`
	ErrorMessage  string     `json:"errorMessage"`
	ProcessDate   *time.Time `json:"processDate"`
	CreatedDate   time.Time  `json:"createdDate"`
	UpdatedDate   *time.Time `json:"updatedDate"`
}
