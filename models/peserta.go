package models

import "gorm.io/gorm"

type Peserta struct {
	gorm.Model
	NamaLengkap    string `gorm:"type:varchar(100);not null" json:"nama_lengkap" binding:"required,min=2,max=100"`
	Email          string `gorm:"uniqueIndex;type:varchar(100);not null" json:"email" binding:"required,email,max=100"`
	Password       string `gorm:"type:varchar(255);not null" json:"-" binding:"required,min=8"`
	JenisPeserta   string `gorm:"type:varchar(50);not null" json:"jenis_peserta" binding:"required,oneof=Mitra Peserta"`
	CabangOlahraga string `gorm:"type:varchar(50);not null" json:"cabang_olahraga" binding:"required"`
	Aset           string `gorm:"type:varchar(100);not null" json:"aset" binding:"required"`
	Foto           string `gorm:"type:varchar(255)" json:"foto"`
	FotoPath       string `json:"foto_path"`
}
