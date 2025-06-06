package controllers

import (
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/models"

	"github.com/gin-gonic/gin"
)

// Handler untuk pendaftaran peserta
func SubmitParticipantRegistration(c *gin.Context) {
	var input models.ParticipantRegistration
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data form tidak valid", "details": err.Error()})
		return
	}

	input.WaktuDaftar = time.Now()

	// Daftar nama field dokumen yang wajib di-upload
	requiredFiles := []string{
		"ktp", "id_card", "bpjs", "pas_foto",
		"surat_pernyataan", "surat_layak_bertanding", "form_prq",
	}

	// Tambahan jika jenis_peserta adalah Mitra
	if strings.ToLower(input.JenisPeserta) == "mitra" {
		requiredFiles = append(requiredFiles, "surat_keterangan_kerja", "kontrak_kerja", "sertifikat_bst")
		if strings.TrimSpace(input.MediaSosial) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Link media sosial wajib diisi untuk peserta mitra"})
			return
		}
	}

	uploadedFiles := make(map[string]string)

	for _, field := range requiredFiles {
		file, err := c.FormFile(field)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Dokumen " + field + " wajib diunggah"})
			return
		}

		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".pdf" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak valid untuk: " + field})
			return
		}

		// Simpan file
		filename := time.Now().Format("20060102150405") + "_" + field + ext
		savePath := "./uploads/" + filename
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file: " + field})
			return
		}

		uploadedFiles[field] = filename
	}

	// Simpan ke DB jika Anda sudah menggunakan GORM
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Pendaftaran berhasil",
		"data":         input,
		"uploadedDocs": uploadedFiles,
	})
}
