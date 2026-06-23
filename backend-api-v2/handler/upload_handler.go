package handler

import (
	"backend-api-v2/service"
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct {
	Svc service.UploadServiceInterface
}

func NewUploadHandler(svc service.UploadServiceInterface) *UploadHandler {
	return &UploadHandler{Svc: svc}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "9999", "message": "กรุณาแนบไฟล์", "data": nil})
		return
	}
	if filepath.Ext(fileHeader.Filename) != ".txt" {
		c.JSON(http.StatusBadRequest, gin.H{"code": "9999", "message": "รองรับเฉพาะไฟล์ .txt", "data": nil})
		return
	}

	file, _ := fileHeader.Open()
	defer file.Close()
	fileBytes, _ := io.ReadAll(file)

	result, err := h.Svc.ProcessUpload(fileHeader.Filename, fileBytes)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": "9999", "message": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0000", "message": "success", "data": result})
}

func (h *UploadHandler) GetHistory(c *gin.Context) {
	fileName := c.Query("fileName")
	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	headers, total, err := h.Svc.GetHistory(fileName, status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "9999", "message": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0000", "message": "success", "data": gin.H{
		"total": total, "page": page, "pageSize": pageSize, "items": headers,
	}})
}

func (h *UploadHandler) GetDetail(c *gin.Context) {
	uploadID, _ := strconv.Atoi(c.Param("uploadId"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	details, total, err := h.Svc.GetDetail(uploadID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": "9999", "message": err.Error(), "data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"code": "0000", "message": "success", "data": gin.H{
		"total": total, "page": page, "pageSize": pageSize, "items": details,
	}})
}
