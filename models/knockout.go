package models

import (
	"time"

	"gorm.io/gorm"
)

type KnockoutMatch struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Kategori    string `gorm:"type:varchar(100);not null" json:"kategori"`
	SubKategori string `gorm:"type:varchar(100);not null" json:"sub_kategori"`
	Tim1        string `gorm:"type:varchar(255);not null" json:"tim1"`
	Tim2        string `gorm:"type:varchar(255);not null" json:"tim2"`
	Hasil       string `gorm:"type:varchar(100);not null" json:"hasil"`
	Tahap       string `gorm:"type:varchar(100);not null" json:"tahap"`
}
