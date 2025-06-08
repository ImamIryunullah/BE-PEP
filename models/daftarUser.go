package models

import (
	"time"

	"gorm.io/gorm"
)

type DaftarUser struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Email     string     `gorm:"uniqueIndex;type:varchar(100);not null" json:"email"`
	Password  string     `gorm:"type:varchar(255);not null" json:"-"`
	Aset      string     `gorm:"type:varchar(100);not null" json:"aset"`
	Provinsi  string     `gorm:"type:varchar(100);not null" json:"provinsi"`
	Foto      string     `gorm:"type:varchar(255)" json:"foto"`
	FotoPath  string     `gorm:"type:varchar(500)" json:"foto_path"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
	LastLogin *time.Time `json:"last_login,omitempty"`

	Registrations []ParticipantRegistration `gorm:"foreignKey:UserID" json:"registrations,omitempty"`
}

func (u *DaftarUser) GetPublicData() map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"email":      u.Email,
		"aset":       u.Aset,
		"provinsi":   u.Provinsi,
		"foto":       u.Foto,
		"is_active":  u.IsActive,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	}
}
