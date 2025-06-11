package seeder

import (
	"time"

	"github.com/ImamIryunullah/BE-PEP/config"

	"github.com/ImamIryunullah/BE-PEP/models"

	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin() {
	password := "admin123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	admin := models.DaftarUser{
		Email:     "admin@example.com",
		Password:  string(hashedPassword),
		Aset:      "Kantor Pusat",
		Provinsi:  "PALEMBANG",
		Role:      "admin",
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	config.DB.FirstOrCreate(&admin, models.DaftarUser{Email: admin.Email})
}
