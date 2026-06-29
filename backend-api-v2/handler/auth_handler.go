package handler

import (
	"backend-api-v2/repository"
	"backend-api-v2/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuditRepo *repository.AuditLogRepository
}

func NewAuthHandler(auditRepo *repository.AuditLogRepository) *AuthHandler {
	return &AuthHandler{AuditRepo: auditRepo}
}

// GET /api/v1/oauth/authorize
func (h *AuthHandler) Authorize(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    "0000",
		"message": "success",
		"data": gin.H{
			"authorizationCode": util.MockAuthCode,
			"expiresAt":         time.Now().Add(5 * time.Minute),
		},
	})
}

// POST /api/v1/oauth/token
func (h *AuthHandler) Token(c *gin.Context) {
	var req struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "9999", "message": "Invalid request", "data": nil,
		})
		return
	}

	if !util.ValidateCode(req.Code) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": "4001", "message": "Invalid authorization code", "data": nil,
		})
		return
	}

	// บันทึก Audit Log
	h.AuditRepo.Insert("intern_user", "LOGIN_SUCCESS", "", "User login successfully")

	c.JSON(http.StatusOK, gin.H{
		"code":    "0000",
		"message": "success",
		"data": gin.H{
			"accessToken": util.MockToken,
			"tokenType":   "Bearer",
			"expiresIn":   3600,
		},
	})
}

// GET /api/v1/me
func (h *AuthHandler) Me(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    "0000",
		"message": "success",
		"data": gin.H{
			"username":    "intern_user",
			"displayName": "Intern User",
			"role":        "user",
		},
	})
}

// POST /api/v1/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	h.AuditRepo.Insert("intern_user", "LOGOUT", "", "User logged out")
	c.JSON(http.StatusOK, gin.H{
		"code": "0000", "message": "success", "data": nil,
	})
}
