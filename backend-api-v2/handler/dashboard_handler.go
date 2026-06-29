package handler

import (
	"backend-api-v2/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	DashboardRepo *repository.DashboardRepository
}

func NewDashboardHandler(dashboardRepo *repository.DashboardRepository) *DashboardHandler {
	return &DashboardHandler{DashboardRepo: dashboardRepo}
}

// GET /api/v1/dashboard/summary
func (h *DashboardHandler) GetSummary(c *gin.Context) {
	summary, err := h.DashboardRepo.GetSummary()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": "9999", "message": err.Error(), "data": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": "0000", "message": "success", "data": summary,
	})
}
