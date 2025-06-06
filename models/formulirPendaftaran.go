package models

import "time"

type ParticipantRegistration struct {
	NamaLengkap    string    `form:"nama_lengkap" binding:"required"`
	Email          string    `form:"email" binding:"required,email"`
	NoTelepon      string    `form:"no_telepon" binding:"required"`
	JenisPeserta   string    `form:"jenis_peserta" binding:"required"`
	CabangOlahraga string    `form:"cabang_olahraga" binding:"required"`
	WilayahKerja   string    `form:"wilayah_kerja" binding:"required"`
	MediaSosial    string    `form:"media_sosial"`
	Catatan        string    `form:"catatan"`
	WaktuDaftar    time.Time `form:"-"`
}
