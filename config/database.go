package config

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
    dsn := "root:root@tcp(127.0.0.1:3306)/pep_db?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("Gagal koneksi ke database: %v", err)
    }

    log.Println("✅ Database connected.")
    DB = db
}
