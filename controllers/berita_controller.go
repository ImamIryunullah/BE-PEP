package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ImamIryunullah/BE-PEP/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type BeritaInput struct {
	Judul    string `json:"judul" binding:"required"`
	Subtitle string `json:"subtitle"`
	Tanggal  string `json:"tanggal" binding:"required"`
	Penulis  string `json:"penulis" binding:"required"`
	Isi      string `json:"isi" binding:"required"`
	Foto     string `json:"foto"`
}

// Updated controller functions
func UpdateBerita(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var berita models.Berita
	if err := db.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	// Parse form data instead of JSON
	judul := c.PostForm("judul")
	subtitle := c.PostForm("subtitle")
	tanggal := c.PostForm("tanggal")
	penulis := c.PostForm("penulis")
	isi := c.PostForm("isi")

	// Validate required fields
	if judul == "" || tanggal == "" || penulis == "" || isi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field yang diperlukan tidak boleh kosong"})
		return
	}

	// Parse date
	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal harus yyyy-mm-dd"})
		return
	}

	// Handle file upload
	foto := berita.Foto // Keep existing photo by default
	file, header, err := c.Request.FormFile("foto")
	if err == nil {
		defer file.Close()

		// Validate file type
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/jpg":  true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		}

		contentType := header.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung"})
			return
		}

		// Validate file size (5MB)
		if header.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file terlalu besar"})
			return
		}

		// Generate unique filename
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), strings.ReplaceAll(header.Filename, ext, ""), ext)

		// Save file
		uploadPath := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(header, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}

		foto = filename

		// Optional: Delete old photo file if it exists and is different
		if berita.Foto != "" && berita.Foto != filename {
			oldPath := filepath.Join("uploads", berita.Foto)
			os.Remove(oldPath) // Ignore error if file doesn't exist
		}
	}

	// Update berita
	updated := models.Berita{
		Judul:    judul,
		Subtitle: subtitle,
		Tanggal:  t,
		Penulis:  penulis,
		Isi:      isi,
		Foto:     foto,
	}

	if err := db.Model(&berita).Updates(updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Fetch updated data to return consistent response
	if err := db.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data yang diperbarui"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    berita,
		"message": "Berita berhasil diperbarui",
	})
}

func CreateBerita(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// Parse form data
	judul := c.PostForm("judul")
	subtitle := c.PostForm("subtitle")
	tanggal := c.PostForm("tanggal")
	penulis := c.PostForm("penulis")
	isi := c.PostForm("isi")

	// Validate required fields
	if judul == "" || tanggal == "" || penulis == "" || isi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field yang diperlukan tidak boleh kosong"})
		return
	}

	// Parse date
	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal harus yyyy-mm-dd"})
		return
	}

	var foto string
	// Handle file upload
	file, header, err := c.Request.FormFile("foto")
	if err == nil {
		defer file.Close()

		// Validate file type
		allowedTypes := map[string]bool{
			"image/jpeg": true,
			"image/jpg":  true,
			"image/png":  true,
			"image/gif":  true,
			"image/webp": true,
		}

		contentType := header.Header.Get("Content-Type")
		if !allowedTypes[contentType] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format file tidak didukung"})
			return
		}

		// Validate file size (5MB)
		if header.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file terlalu besar"})
			return
		}

		// Generate unique filename
		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), strings.ReplaceAll(header.Filename, ext, ""), ext)

		// Save file
		uploadPath := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(header, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}

		foto = filename
	}

	// Create berita
	berita := models.Berita{
		Judul:    judul,
		Subtitle: subtitle,
		Tanggal:  t,
		Penulis:  penulis,
		Isi:      isi,
		Foto:     foto,
	}

	if err := db.Create(&berita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    berita,
		"message": "Berita berhasil disimpan",
	})
}

func GetAllBerita(c *gin.Context) {
	var beritas []models.Berita
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Order("tanggal desc").Find(&beritas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data berita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": beritas})
}

func DeleteBerita(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var berita models.Berita
	if err := db.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	// Delete associated photo file if exists
	if berita.Foto != "" {
		photoPath := filepath.Join("uploads", berita.Foto)
		os.Remove(photoPath) // Ignore error if file doesn't exist
	}

	if err := db.Delete(&berita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berita berhasil dihapus"})
}
