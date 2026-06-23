package model

import "time"

type UploadHeader struct {
	UploadID     int        `json:"uploadId"`
	FileName     string     `json:"fileName"`
	UploadDate   time.Time  `json:"uploadDate"`
	TotalRecord  int        `json:"totalRecord"`
	TotalSuccess int        `json:"totalSuccess"`
	TotalFailed  int        `json:"totalFailed"`
	Status       string     `json:"status"`
	CreatedDate  time.Time  `json:"createdDate"`
	UpdatedDate  *time.Time `json:"updatedDate"`
}
