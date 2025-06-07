package controllers

import (
	"log"
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
		log.Printf("Binding error: %v", err)
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
			log.Printf("File save error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file: " + field})
			return
		}

		uploadedFiles[field] = filename

		// Set filename ke struct berdasarkan field
		setDocumentField(&input, field, filename)
	}

	// Gunakan konsisten database connection
	if err := config.DB.Create(&input).Error; err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data ke database", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Pendaftaran berhasil",
		"data":         input,
		"uploadedDocs": uploadedFiles,
	})
}

// Helper function untuk set document field
func setDocumentField(input *models.ParticipantRegistration, field, filename string) {
	switch field {
	case "ktp":
		input.KTP = filename
	case "id_card":
		input.IDCard = filename
	case "bpjs":
		input.BPJS = filename
	case "pas_foto":
		input.PasFoto = filename
	case "surat_pernyataan":
		input.SuratPernyataan = filename
	case "surat_layak_bertanding":
		input.SuratLayakBertanding = filename
	case "form_prq":
		input.FormPRQ = filename
	case "surat_keterangan_kerja":
		input.SuratKeteranganKerja = filename
	case "kontrak_kerja":
		input.KontrakKerja = filename
	case "sertifikat_bst":
		input.SertifikatBST = filename
	}
}

func GetAllPeserta(c *gin.Context) {
	// Gunakan konsisten database connection
	var peserta []models.ParticipantRegistration

	if err := config.DB.Order("created_at desc").Find(&peserta).Error; err != nil {
		log.Printf("Database query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal mengambil data peserta",
			"details": err.Error(),
		})
		return
	}

	// Response yang konsisten dengan frontend
	c.JSON(http.StatusOK, gin.H{
		"message": "Data peserta berhasil diambil",
		"data":    peserta,
		"total":   len(peserta),
	})
}
