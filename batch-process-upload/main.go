package main

import (
	"batch-process-upload/repository"
	"batch-process-upload/service"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type Config struct {
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
	log.Println("=== Batch Process Upload เริ่มทำงาน ===")

	f, err := os.ReadFile("configs/config.yaml")
	if err != nil {
		log.Fatal("โหลด config ไม่ได้:", err)
	}
	var cfg Config
	yaml.Unmarshal(f, &cfg)

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

	repo := repository.NewBatchRepository(db, cfg.Database.Schema)
	svc := service.NewBatchService(repo)

	if err := svc.ProcessAll(); err != nil {
		log.Fatal("Batch error:", err)
	}

	log.Println("=== Batch Process Upload เสร็จสิ้น ===")
}
