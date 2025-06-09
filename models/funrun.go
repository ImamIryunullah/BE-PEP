package models

import (
	"time"

	"gorm.io/gorm"
)

type Funrun struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Nama      string `gorm:"type:varchar(100);not null" json:"nama"`
	Email     string `gorm:"type:varchar(100);unique;not null" json:"email"`
	Kontingen string `gorm:"type:varchar(100);not null" json:"kontingen"`
	Size      string `gorm:"type:varchar(10);not null" json:"size"`
	Status    string `json:"status" gorm:"default:pending"`
}
