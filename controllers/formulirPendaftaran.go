package controllers

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/models"

	"github.com/gin-gonic/gin"
)

type ParticipantRegistrationRequest struct {
	NamaLengkap    string `json:"nama_lengkap" form:"nama_lengkap" binding:"required"`
	Email          string `json:"email" form:"email" binding:"required,email"`
	NoTelepon      string `json:"no_telepon" form:"no_telepon" binding:"required"`
	JenisKelamin   string `json:"jenis_kelamin" form:"jenis_kelamin" binding:"required"`
	JenisPeserta   string `json:"jenis_peserta" form:"jenis_peserta" binding:"required"`
	CabangOlahraga string `json:"cabang_olahraga" form:"cabang_olahraga" binding:"required"`
	WilayahKerja   string `json:"wilayah_kerja" form:"wilayah_kerja" binding:"required"`
	MediaSosial    string `json:"media_sosial" form:"media_sosial"`
	Catatan        string `json:"catatan" form:"catatan"`
}

func SubmitParticipantRegistration(c *gin.Context) {
	var input ParticipantRegistrationRequest
	if err := c.ShouldBind(&input); err != nil {
		log.Printf("Binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data form tidak valid", "details": err.Error()})
		return
	}

	var user models.DaftarUser
	if err := config.DB.Where("email = ?", strings.ToLower(input.Email)).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email tidak terdaftar. Silakan daftar akun terlebih dahulu"})
		return
	}

	registration := models.ParticipantRegistration{
		UserID:      user.ID,
		NamaLengkap: input.NamaLengkap,

		NoTelepon:      input.NoTelepon,
		JenisKelamin:   input.JenisKelamin,
		JenisPeserta:   input.JenisPeserta,
		CabangOlahraga: input.CabangOlahraga,
		WilayahKerja:   input.WilayahKerja,
		MediaSosial:    input.MediaSosial,
		Catatan:        input.Catatan,
		WaktuDaftar:    time.Now(),
	}

	requiredFiles := []string{
		"ktp", "id_card", "bpjs", "pas_foto",
		"surat_pernyataan", "surat_layak_bertanding", "form_prq",
	}

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

		filename := time.Now().Format("20060102150405") + "_" + field + ext
		savePath := "./uploads/" + filename
		if err := c.SaveUploadedFile(file, savePath); err != nil {
			log.Printf("File save error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file: " + field})
			return
		}

		uploadedFiles[field] = filename

		setDocumentField(&registration, field, filename)
	}

	if err := config.DB.Create(&registration).Error; err != nil {
		log.Printf("Database error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data ke database", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Pendaftaran berhasil",
		"data":         registration,
		"uploadedDocs": uploadedFiles,
	})
}

func EditParticipantRegistration(c *gin.Context) {
	registrationIDStr := c.Param("id")
	registrationID, err := strconv.ParseUint(registrationIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID registrasi tidak valid"})
		return
	}

	var existing models.ParticipantRegistration
	if err := config.DB.First(&existing, uint(registrationID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data registrasi tidak ditemukan"})
		return
	}

	var input ParticipantRegistrationRequest
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data form tidak valid", "details": err.Error()})
		return
	}

	existing.NamaLengkap = input.NamaLengkap
	existing.NoTelepon = input.NoTelepon
	existing.JenisKelamin = input.JenisKelamin
	existing.JenisPeserta = input.JenisPeserta
	existing.CabangOlahraga = input.CabangOlahraga
	existing.WilayahKerja = input.WilayahKerja
	existing.MediaSosial = input.MediaSosial
	existing.Catatan = input.Catatan

	requiredFiles := []string{
		"ktp", "id_card", "bpjs", "pas_foto",
		"surat_pernyataan", "surat_layak_bertanding", "form_prq",
	}

	if strings.ToLower(input.JenisPeserta) == "mitra" {
		requiredFiles = append(requiredFiles, "surat_keterangan_kerja", "kontrak_kerja", "sertifikat_bst")
		if strings.TrimSpace(input.MediaSosial) == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Link media sosial wajib diisi untuk peserta mitra"})
			return
		}
	}

	updatedFiles := make(map[string]string)

	for _, field := range requiredFiles {
		file, err := c.FormFile(field)
		if err == nil {
			ext := strings.ToLower(filepath.Ext(file.Filename))
			if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".pdf" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak valid untuk: " + field})
				return
			}

			filename := time.Now().Format("20060102150405") + "_" + field + ext
			savePath := "./uploads/" + filename
			if err := c.SaveUploadedFile(file, savePath); err != nil {
				log.Printf("File save error: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file: " + field})
				return
			}

			setDocumentField(&existing, field, filename)
			updatedFiles[field] = filename
		}
	}

	if err := config.DB.Save(&existing).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui data registrasi", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "Data registrasi berhasil diperbarui",
		"data":         existing,
		"updatedFiles": updatedFiles,
	})
}

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
	var peserta []models.ParticipantRegistration

	if err := config.DB.Preload("User").Order("created_at desc").Find(&peserta).Error; err != nil {
		log.Printf("Database query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal mengambil data peserta",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data peserta berhasil diambil",
		"data":    peserta,
		"total":   len(peserta),
	})
}

func GetRegistrationsByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var registrations []models.ParticipantRegistration
	if err := config.DB.Where("user_id = ?", uint(userID)).Order("created_at desc").Find(&registrations).Error; err != nil {
		log.Printf("Database query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal mengambil data registrasi",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data registrasi berhasil diambil",
		"data":    registrations,
		"total":   len(registrations),
	})
}

func GetUserWithRegistrations(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.DaftarUser
	if err := config.DB.Preload("Registrations").First(&user, uint(userID)).Error; err != nil {
		log.Printf("Database query error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Gagal mengambil data user",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data user dengan registrasi berhasil diambil",
		"data":    user.GetPublicDataWithRegistrations(),
	})
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
	Reason string `json:"reason,omitempty"`
}

func UpdateParticipantStatus(c *gin.Context) {
	participantIDStr := c.Param("id")
	participantID, err := strconv.ParseUint(participantIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID peserta tidak valid"})
		return
	}

	var request UpdateStatusRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data request tidak valid", "details": err.Error()})
		return
	}

	allowedStatus := []string{"pending", "approved", "rejected"}
	isValidStatus := false
	for _, status := range allowedStatus {
		if request.Status == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid. Gunakan: pending, approved, atau rejected"})
		return
	}

	var participant models.ParticipantRegistration
	if err := config.DB.First(&participant, uint(participantID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data peserta tidak ditemukan"})
		return
	}

	participant.Status = request.Status

	if request.Status == "rejected" && request.Reason != "" {
		participant.Catatan = request.Reason
	}

	if err := config.DB.Save(&participant).Error; err != nil {
		log.Printf("Database update error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status peserta", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status peserta berhasil diupdate",
		"data": gin.H{
			"id":     participant.ID,
			"status": participant.Status,
			"nama":   participant.NamaLengkap,
		},
	})
}

func GetParticipantById(c *gin.Context) {
	participantIDStr := c.Param("id")
	participantID, err := strconv.ParseUint(participantIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID peserta tidak valid"})
		return
	}

	var participant models.ParticipantRegistration
	if err := config.DB.Preload("User").First(&participant, uint(participantID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data peserta tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Data peserta berhasil diambil",
		"data":    participant,
	})
}
