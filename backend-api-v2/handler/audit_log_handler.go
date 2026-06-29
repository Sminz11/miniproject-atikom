package handler

import (
	"backend-api-v2/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuditLogHandler struct {
	AuditRepo *repository.AuditLogRepository
}

func NewAuditLogHandler(auditRepo *repository.AuditLogRepository) *AuditLogHandler {
	return &AuditLogHandler{AuditRepo: auditRepo}
}

// GET /api/v1/audit-logs
func (h *AuditLogHandler) GetLogs(c *gin.Context) {
	action := c.Query("action")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	logs, total, err := h.AuditRepo.GetLogs(action, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "9999", "message": err.Error(), "data": nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    "0000",
		"message": "success",
		"data": gin.H{
			"total": total, "page": page, "pageSize": pageSize, "items": logs,
		},
	})
}
