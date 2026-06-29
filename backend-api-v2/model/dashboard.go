package model

type DashboardSummary struct {
	TotalUpload    int            `json:"totalUpload"`
	TotalSuccess   int            `json:"totalSuccess"`
	TotalFailed    int            `json:"totalFailed"`
	TotalPending   int            `json:"totalPending"`
	UploadByStatus map[string]int `json:"uploadByStatus"`
	RecentUploads  []UploadHeader `json:"recentUploads"`
}
