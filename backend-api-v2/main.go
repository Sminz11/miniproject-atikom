package main

import (
	"backend-api-v2/handler"
	"backend-api-v2/middleware"
	"backend-api-v2/repository"
	"backend-api-v2/service"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Port int `yaml:"port"`
	} `yaml:"app"`
	Database struct {
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Schema   string `yaml:"schema"`
	} `yaml:"database"`
}

func main() {
	// โหลด Config
	f, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatal("โหลด config ไม่ได้:", err)
	}
	var cfg Config
	yaml.Unmarshal(f, &cfg)

	// เชื่อมต่อ Database
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username,
		cfg.Database.Password, cfg.Database.DBName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("เชื่อมต่อ DB ไม่ได้:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Ping DB ไม่ได้:", err)
	}
	log.Println("เชื่อมต่อ Database สำเร็จ")

	// สร้าง Repository
	uploadRepo := repository.NewUploadRepository(db, cfg.Database.Schema)
	auditRepo := repository.NewAuditLogRepository(db, cfg.Database.Schema)
	dashboardRepo := repository.NewDashboardRepository(db, cfg.Database.Schema)
	retryRepo := repository.NewRetryRepository(db, cfg.Database.Schema)

	// สร้าง Service
	uploadSvc := service.NewUploadService(uploadRepo)

	// สร้าง Handler
	uploadHandler := handler.NewUploadHandler(uploadSvc)
	authHandler := handler.NewAuthHandler(auditRepo)
	dashboardHandler := handler.NewDashboardHandler(dashboardRepo)
	retryHandler := handler.NewRetryHandler(retryRepo, auditRepo)
	auditLogHandler := handler.NewAuditLogHandler(auditRepo)

	// สร้าง Router
	r := gin.Default()

	// CORS Middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api/v1")
	{
		// OAuth Routes (ไม่ต้อง Auth)
		oauth := api.Group("/oauth")
		{
			oauth.GET("/authorize", authHandler.Authorize)
			oauth.POST("/token", authHandler.Token)
		}

		// Protected Routes (ต้อง Auth)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
			protected.GET("/me", authHandler.Me)
			protected.POST("/logout", authHandler.Logout)

			// Upload Routes
			protected.POST("/uploads", uploadHandler.Upload)
			protected.GET("/uploads", uploadHandler.GetHistory)
			protected.GET("/uploads/:uploadId/details", uploadHandler.GetDetail)
			protected.POST("/uploads/:uploadId/details/:detailId/retry", retryHandler.Retry)
			protected.GET("/uploads/:uploadId/details/export", retryHandler.ExportCSV)

			// Dashboard
			protected.GET("/dashboard/summary", dashboardHandler.GetSummary)

			// Audit Log
			protected.GET("/audit-logs", auditLogHandler.GetLogs)
		}
	}

	log.Printf("Server เริ่มที่ port %d", cfg.App.Port)
	r.Run(fmt.Sprintf(":%d", cfg.App.Port))
}
