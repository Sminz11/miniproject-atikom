package main

import (
	"backend-api-v2/handler"
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

	// สร้าง Layer ต่างๆ
	repo := repository.NewUploadRepository(db, cfg.Database.Schema)
	svc := service.NewUploadService(repo)
	h := handler.NewUploadHandler(svc)

	// สร้าง Router
	r := gin.Default()

	// CORS Middleware - อนุญาตให้ Angular (localhost:4200) เรียก API ได้
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
		api.POST("/uploads", h.Upload)
		api.GET("/uploads", h.GetHistory)
		api.GET("/uploads/:uploadId/details", h.GetDetail)
	}

	log.Printf("Server เริ่มที่ port %d", cfg.App.Port)
	r.Run(fmt.Sprintf(":%d", cfg.App.Port))
}
