package models

import (
	"time"

	"gorm.io/gorm"
)

type ParticipantRegistration struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         uint      `json:"user_id" gorm:"not null;index"`
	NamaLengkap    string    `json:"nama_lengkap" form:"nama_lengkap" binding:"required"`
	Email          string    `json:"email" form:"email" binding:"required,email"`
	NoTelepon      string    `json:"no_telepon" form:"no_telepon" binding:"required"`
	JenisPeserta   string    `json:"jenis_peserta" form:"jenis_peserta" binding:"required"`
	CabangOlahraga string    `json:"cabang_olahraga" form:"cabang_olahraga" binding:"required"`
	WilayahKerja   string    `json:"wilayah_kerja" form:"wilayah_kerja" binding:"required"`
	MediaSosial    string    `json:"media_sosial" form:"media_sosial"`
	Catatan        string    `json:"catatan" form:"catatan"`
	WaktuDaftar    time.Time `json:"waktu_daftar" form:"-"`

	KTP                  string `json:"ktp,omitempty"`
	IDCard               string `json:"id_card,omitempty"`
	BPJS                 string `json:"bpjs,omitempty"`
	PasFoto              string `json:"pas_foto,omitempty"`
	SuratPernyataan      string `json:"surat_pernyataan,omitempty"`
	SuratLayakBertanding string `json:"surat_layak_bertanding,omitempty"`
	FormPRQ              string `json:"form_prq,omitempty"`
	SuratKeteranganKerja string `json:"surat_keterangan_kerja,omitempty"`
	KontrakKerja         string `json:"kontrak_kerja,omitempty"`
	SertifikatBST        string `json:"sertifikat_bst,omitempty"`

	Status string `json:"status" gorm:"default:pending"`

	User DaftarUser `gorm:"foreignKey:UserID" json:"user,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}

func (u *DaftarUser) GetUserRegistrations(db *gorm.DB) ([]ParticipantRegistration, error) {
	var registrations []ParticipantRegistration
	err := db.Where("user_id = ?", u.ID).Find(&registrations).Error
	return registrations, err
}

func (u *DaftarUser) GetPublicDataWithRegistrations() map[string]interface{} {
	data := u.GetPublicData()
	data["registrations"] = u.Registrations
	return data
}
