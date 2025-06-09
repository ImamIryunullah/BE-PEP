package controllers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/models"
	"github.com/gin-gonic/gin"
)

type BeritaInput struct {
	Judul    string `json:"judul" binding:"required"`
	Subtitle string `json:"subtitle"`
	Tanggal  string `json:"tanggal" binding:"required"`
	Penulis  string `json:"penulis" binding:"required"`
	Isi      string `json:"isi" binding:"required"`
	Foto     string `json:"foto"`
}

func UpdateBerita(c *gin.Context) {
	id := c.Param("id")
	var berita models.Berita
	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	judul := c.PostForm("judul")
	subtitle := c.PostForm("subtitle")
	tanggal := c.PostForm("tanggal")
	penulis := c.PostForm("penulis")
	isi := c.PostForm("isi")

	if judul == "" || tanggal == "" || penulis == "" || isi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field yang diperlukan tidak boleh kosong"})
		return
	}

	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal harus yyyy-mm-dd"})
		return
	}

	foto := berita.Foto
	file, header, err := c.Request.FormFile("foto")
	if err == nil {
		defer file.Close()

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

		if header.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file terlalu besar"})
			return
		}

		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), strings.ReplaceAll(header.Filename, ext, ""), ext)

		uploadPath := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(header, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}

		foto = filename

		if berita.Foto != "" && berita.Foto != filename {
			oldPath := filepath.Join("uploads", berita.Foto)
			os.Remove(oldPath)
		}
	}

	updated := models.Berita{
		Judul:    judul,
		Subtitle: subtitle,
		Tanggal:  t,
		Penulis:  penulis,
		Isi:      isi,
		Foto:     foto,
	}

	if err := config.DB.Model(&berita).Updates(updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data yang diperbarui"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":    berita,
		"message": "Berita berhasil diperbarui",
	})
}

func CreateBerita(c *gin.Context) {
	judul := c.PostForm("judul")
	subtitle := c.PostForm("subtitle")
	tanggal := c.PostForm("tanggal")
	penulis := c.PostForm("penulis")
	isi := c.PostForm("isi")

	if judul == "" || tanggal == "" || penulis == "" || isi == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Field yang diperlukan tidak boleh kosong"})
		return
	}

	t, err := time.Parse("2006-01-02", tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal harus yyyy-mm-dd"})
		return
	}

	var foto string
	file, header, err := c.Request.FormFile("foto")
	if err == nil {
		defer file.Close()

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

		if header.Size > 5*1024*1024 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Ukuran file terlalu besar"})
			return
		}

		ext := filepath.Ext(header.Filename)
		filename := fmt.Sprintf("%d_%s%s", time.Now().Unix(), strings.ReplaceAll(header.Filename, ext, ""), ext)

		uploadPath := filepath.Join("uploads", filename)
		if err := c.SaveUploadedFile(header, uploadPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan file"})
			return
		}

		foto = filename
	}

	berita := models.Berita{
		Judul:    judul,
		Subtitle: subtitle,
		Tanggal:  t,
		Penulis:  penulis,
		Isi:      isi,
		Foto:     foto,
	}

	if err := config.DB.Create(&berita).Error; err != nil {
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

	if err := config.DB.Order("tanggal desc").Find(&beritas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data berita"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": beritas})
}
func GetBeritaByID(c *gin.Context) {
	id := c.Param("id")
	var berita models.Berita

	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": berita})
}
func DeleteBerita(c *gin.Context) {

	id := c.Param("id")

	var berita models.Berita
	if err := config.DB.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	if berita.Foto != "" {
		photoPath := filepath.Join("uploads", berita.Foto)
		os.Remove(photoPath)
	}

	if err := config.DB.Delete(&berita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berita berhasil dihapus"})
}
