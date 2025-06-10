package models

import (
	"time"

	"gorm.io/gorm"
)

type Berita struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Judul     string         `gorm:"type:varchar(255);not null" json:"judul"`
	Subtitle  string         `gorm:"type:varchar(255)" json:"subtitle"`
	Tanggal   time.Time      `json:"tanggal"`
	Penulis   string         `gorm:"type:varchar(100);not null" json:"penulis"`
	Isi       string         `gorm:"type:text;not null" json:"isi"`
	Foto      string         `gorm:"type:varchar(255)" json:"foto"`
}
