package controllers

import (
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegistrationRequest struct {
	Email    string                `form:"email" binding:"required,email,max=100"`
	Password string                `form:"password" binding:"required,min=8"`
	Aset     string                `form:"aset" binding:"required,max=100"`
	Provinsi string                `form:"provinsi" binding:"required,max=100"`
	Foto     *multipart.FileHeader `form:"foto" binding:"required"`
}

func RegisterPeserta(c *gin.Context) {
	var req RegistrationRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	if err := validatePhoto(req.Foto); err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	uploadDir := "uploads"
	if err := ensureDirectoryExists(uploadDir); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Gagal menyiapkan direktori upload",
		})
		return
	}

	filename := generateUniqueFilename(req.Foto.Filename)
	filepath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(req.Foto, filepath); err != nil {
		c.JSON(500, gin.H{
			"success": false,
			"message": "Gagal menyimpan file foto",
		})
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {

		os.Remove(filepath)
		c.JSON(500, gin.H{
			"success": false,
			"message": "Gagal mengenkripsi password",
		})
		return
	}

	peserta := models.DaftarUser{
		Email:    strings.ToLower(strings.TrimSpace(req.Email)),
		Password: hashedPassword,
		Aset:     strings.TrimSpace(req.Aset),
		Provinsi: strings.TrimSpace(req.Provinsi),
		Foto:     filename,
		FotoPath: filepath,
	}

	if err := savePeserta(&peserta); err != nil {

		os.Remove(filepath)

		if isDuplicateError(err) {
			c.JSON(409, gin.H{
				"success": false,
				"message": "Email sudah terdaftar",
			})
			return
		}

		c.JSON(500, gin.H{
			"success": false,
			"message": "Gagal menyimpan data pendaftaran",
		})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"message": "Pendaftaran berhasil",
		"data": gin.H{
			"id":       peserta.ID,
			"email":    peserta.Email,
			"aset":     peserta.Aset,
			"provinsi": peserta.Provinsi,
			"foto":     peserta.Foto,
		},
	})
}

func validatePhoto(file *multipart.FileHeader) error {

	maxSize := int64(5 * 1024 * 1024)
	if file.Size > maxSize {
		return fmt.Errorf("ukuran file terlalu besar, maksimal 5MB")
	}

	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return fmt.Errorf("format file tidak didukung. Gunakan: jpg, jpeg, png, gif, atau webp")
	}

	// Validasi MIME type
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}

	// Buka file untuk cek header
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("gagal membaca file")
	}
	defer src.Close()

	// Baca 512 bytes pertama untuk deteksi MIME type
	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil {
		return fmt.Errorf("gagal membaca header file")
	}

	// Reset file pointer
	src.Seek(0, 0)

	// Deteksi MIME type
	contentType := DetectContentType(buffer)
	if !allowedMimeTypes[contentType] {
		return fmt.Errorf("tipe file tidak valid")
	}

	return nil
}

// DetectContentType mendeteksi content type dari buffer
func DetectContentType(buffer []byte) string {
	// Implementasi sederhana deteksi content type
	if len(buffer) >= 3 && buffer[0] == 0xFF && buffer[1] == 0xD8 && buffer[2] == 0xFF {
		return "image/jpeg"
	}
	if len(buffer) >= 8 && string(buffer[0:8]) == "\x89PNG\r\n\x1a\n" {
		return "image/png"
	}
	if len(buffer) >= 6 && string(buffer[0:6]) == "GIF87a" || string(buffer[0:6]) == "GIF89a" {
		return "image/gif"
	}
	if len(buffer) >= 12 && string(buffer[8:12]) == "WEBP" {
		return "image/webp"
	}
	return "application/octet-stream"
}

// ensureDirectoryExists memastikan direktori ada
func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

// generateUniqueFilename membuat nama file yang unik
func generateUniqueFilename(originalName string) string {
	ext := filepath.Ext(originalName)
	nameWithoutExt := strings.TrimSuffix(originalName, ext)

	// Sanitize filename
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, " ", "_")
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, "..", "")

	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d_%s%s", timestamp, nameWithoutExt, ext)
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func savePeserta(peserta *models.DaftarUser) error {
	return config.DB.Create(peserta).Error
}

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}

	errorStr := strings.ToLower(err.Error())
	return strings.Contains(errorStr, "duplicate") ||
		strings.Contains(errorStr, "unique") ||
		strings.Contains(errorStr, "constraint")
}

func GetPesertaByEmail(email string) (*models.DaftarUser, error) {
	var peserta models.DaftarUser
	err := config.DB.Where("email = ?", strings.ToLower(strings.TrimSpace(email))).First(&peserta).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &peserta, nil
}

func UpdatePeserta(peserta *models.DaftarUser) error {
	return config.DB.Save(peserta).Error
}

func DeletePeserta(id uint) error {
	var peserta models.DaftarUser
	if err := config.DB.First(&peserta, id).Error; err != nil {
		return err
	}

	if peserta.FotoPath != "" {
		os.Remove(peserta.FotoPath)
	}

	return config.DB.Delete(&peserta).Error
}
