package controllers

import (
	"net/http"
	"time"

	"github.com/ImamIryunullah/BE-PEP/config"
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

func GetBeritas(c *gin.Context) {
	var beritas []models.Berita
	db := c.MustGet("db").(*gorm.DB)
	db.Order("tanggal desc").Find(&beritas)
	c.JSON(http.StatusOK, gin.H{"data": beritas})
}

func CreateBerita(c *gin.Context) {
	judul := c.PostForm("judul")
	subtitle := c.PostForm("subtitle")
	tanggalStr := c.PostForm("tanggal")
	penulis := c.PostForm("penulis")
	isi := c.PostForm("isi")

	tanggal, err := time.Parse("2006-01-02", tanggalStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid. Gunakan format YYYY-MM-DD"})
		return
	}

	file, err := c.FormFile("foto")
	var filename string
	if err == nil {
		filename = file.Filename
		c.SaveUploadedFile(file, "./uploads/"+filename)
	}

	berita := models.Berita{
		Judul:    judul,
		Subtitle: subtitle,
		Tanggal:  tanggal,
		Penulis:  penulis,
		Isi:      isi,
		Foto:     filename,
	}

	if err := config.DB.Create(&berita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan berita ke database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berita berhasil disimpan", "data": berita})
}
func UpdateBerita(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	var berita models.Berita
	if err := db.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	var input BeritaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	t, err := time.Parse("2006-01-02", input.Tanggal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal harus yyyy-mm-dd"})
		return
	}

	updated := models.Berita{
		Judul:    input.Judul,
		Subtitle: input.Subtitle,
		Tanggal:  t,
		Penulis:  input.Penulis,
		Isi:      input.Isi,
		Foto:     input.Foto,
	}

	if err := db.Model(&berita).Updates(updated).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": berita})
}
func GetAllBerita(c *gin.Context) {

	beritaList := []models.Berita{
		{
			ID:       1,
			Judul:    "Peluncuran Twibbon Mini Olympic 2025",
			Subtitle: "Berita Twibbon",
			Tanggal:  time.Now(),
			Penulis:  "Admin PEP",
			Isi:      "Mini Olympic 2025 akan segera dimulai...",
			Foto:     "foto1.jpg",
		},
	}

	c.JSON(http.StatusOK, gin.H{"data": beritaList})
}

func DeleteBerita(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	var berita models.Berita
	if err := db.First(&berita, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Berita tidak ditemukan"})
		return
	}

	if err := db.Delete(&berita).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": true})
}
