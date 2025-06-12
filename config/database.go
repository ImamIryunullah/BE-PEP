package config

import (
	"log"

	"github.com/ImamIryunullah/BE-PEP/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "root:password@tcp(127.0.0.1:3306)/pep_db?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi ke database: %v", err)
	}
	if err := db.AutoMigrate(
		&models.Berita{},
		&models.DaftarUser{},
		&models.ParticipantRegistration{},
		&models.Funrun{},
		&models.SessionLogin{},
		&models.KnockoutMatch{},
	); err != nil {
		log.Fatalf("Gagal melakukan migrasi: %v", err)
	}

	log.Println("Migrasi database berhasil!")
	log.Println("✅ Database connected.")
	DB = db
}
