package controllers

import (
	"fmt"
	"strings"
	"time"

	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterPeserta(c *gin.Context) {
	var peserta models.Peserta

	// Ambil data yang sesuai dengan model
	peserta.NamaLengkap = c.PostForm("nama_lengkap")
	peserta.Email = c.PostForm("email")
	peserta.Password = c.PostForm("password")
	peserta.JenisPeserta = c.PostForm("jenis_peserta")
	peserta.CabangOlahraga = c.PostForm("cabang_olahraga")
	peserta.Aset = c.PostForm("aset")

	// Validasi manual awal
	if peserta.NamaLengkap == "" || peserta.Email == "" || peserta.Password == "" || 
		peserta.JenisPeserta == "" || peserta.CabangOlahraga == "" || peserta.Aset == "" {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Semua field wajib diisi",
		})
		return
	}

	file, err := c.FormFile("foto")
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "File foto harus diupload",
		})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filepath := "uploads/" + filename

	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Gagal menyimpan file",
		})
		return
	}

	peserta.Foto = filename
	peserta.FotoPath = filepath // tambahkan jika kamu ingin menyimpan path relatif/absolut

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(peserta.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Gagal mengenkripsi password",
		})
		return
	}
	peserta.Password = string(hashedPassword)

	if err := config.DB.Create(&peserta).Error; err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			c.JSON(409, gin.H{
				"success": false,
				"message": "Email sudah terdaftar",
			})
			return
		}
		c.JSON(500, gin.H{
			"success": false,
			"message": "Gagal menyimpan data",
		})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "Pendaftaran berhasil",
		"data": gin.H{
			"id":             peserta.ID,
			"nama_lengkap":   peserta.NamaLengkap,
			"email":          peserta.Email,
			"jenis_peserta":  peserta.JenisPeserta,
			"cabang_olahraga": peserta.CabangOlahraga,
			"aset":           peserta.Aset,
			"foto":           peserta.Foto,
		},
	})
}
