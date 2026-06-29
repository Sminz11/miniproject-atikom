package handler

import (
	"backend-api-v2/repository"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RetryHandler struct {
	RetryRepo *repository.RetryRepository
	AuditRepo *repository.AuditLogRepository
}

func NewRetryHandler(retryRepo *repository.RetryRepository, auditRepo *repository.AuditLogRepository) *RetryHandler {
	return &RetryHandler{RetryRepo: retryRepo, AuditRepo: auditRepo}
}

// POST /api/v1/uploads/:uploadId/details/:detailId/retry
func (h *RetryHandler) Retry(c *gin.Context) {
	uploadID, _ := strconv.Atoi(c.Param("uploadId"))
	detailID, _ := strconv.Atoi(c.Param("detailId"))

	if err := h.RetryRepo.RetryDetail(uploadID, detailID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "9999", "message": err.Error(), "data": nil,
		})
		return
	}

	username, _ := c.Get("username")
	h.AuditRepo.Insert(
		fmt.Sprintf("%v", username),
		"RETRY_TRANSACTION",
		fmt.Sprintf("%d", detailID),
		fmt.Sprintf("Retry detail ID %d of upload ID %d", detailID, uploadID),
	)

	c.JSON(http.StatusOK, gin.H{
		"code": "0000", "message": "Retry สำเร็จ รอ Batch ประมวลผล", "data": nil,
	})
}

// GET /api/v1/uploads/:uploadId/details/export
func (h *RetryHandler) ExportCSV(c *gin.Context) {
	uploadID, _ := strconv.Atoi(c.Param("uploadId"))

	results, err := h.RetryRepo.GetDetailsForExport(uploadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "9999", "message": err.Error(), "data": nil,
		})
		return
	}

	// สร้าง CSV
	csv := "detail_id,transaction_id,amount,mobile_no,status,error_message,process_date,created_date\n"
	for _, r := range results {
		processDate := ""
		if r["processDate"] != nil {
			processDate = fmt.Sprintf("%v", r["processDate"])
		}
		csv += fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v\n",
			r["detailId"], r["transactionId"], r["amount"],
			r["mobileNo"], r["status"], r["errorMessage"],
			processDate, r["createdDate"])
	}

	username, _ := c.Get("username")
	h.AuditRepo.Insert(
		fmt.Sprintf("%v", username),
		"EXPORT_CSV",
		fmt.Sprintf("%d", uploadID),
		fmt.Sprintf("Export CSV for upload ID %d", uploadID),
	)

	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=upload_%d.csv", uploadID))
	c.String(http.StatusOK, csv)
}
