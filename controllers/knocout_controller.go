package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CreateKnockoutMatchRequest struct {
	Kategori    string `json:"kategori" binding:"required,max=100"`
	SubKategori string `json:"sub_kategori" binding:"required,max=100"`
	Tim1        string `json:"tim1" binding:"required,max=255"`
	Tim2        string `json:"tim2" binding:"required,max=255"`
	Hasil       string `json:"hasil" binding:"required,max=100"`
	Tahap       string `json:"tahap" binding:"required,max=100"`
}

type UpdateKnockoutMatchRequest struct {
	Kategori    string `json:"kategori" binding:"max=100"`
	SubKategori string `json:"sub_kategori" binding:"max=100"`
	Tim1        string `json:"tim1" binding:"max=255"`
	Tim2        string `json:"tim2" binding:"max=255"`
	Hasil       string `json:"hasil" binding:"max=100"`
	Tahap       string `json:"tahap" binding:"max=100"`
}

func CreateKnockoutMatch(c *gin.Context) {
	var req CreateKnockoutMatchRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	req.Kategori = strings.TrimSpace(req.Kategori)
	req.SubKategori = strings.TrimSpace(req.SubKategori)
	req.Tim1 = strings.TrimSpace(req.Tim1)
	req.Tim2 = strings.TrimSpace(req.Tim2)
	req.Hasil = strings.TrimSpace(req.Hasil)
	req.Tahap = strings.TrimSpace(req.Tahap)

	knockoutMatch := models.KnockoutMatch{
		Kategori:    req.Kategori,
		SubKategori: req.SubKategori,
		Tim1:        req.Tim1,
		Tim2:        req.Tim2,
		Hasil:       req.Hasil,
		Tahap:       req.Tahap,
	}

	if err := config.DB.Create(&knockoutMatch).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menyimpan data pertandingan knockout",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Pertandingan knockout berhasil ditambahkan",
		"data":    knockoutMatch,
	})
}

func GetAllKnockoutMatches(c *gin.Context) {
	var matches []models.KnockoutMatch
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	kategori := c.Query("kategori")
	subKategori := c.Query("sub_kategori")
	tahap := c.Query("tahap")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := config.DB.Model(&models.KnockoutMatch{})

	if kategori != "" {
		query = query.Where("kategori ILIKE ?", "%"+kategori+"%")
	}
	if subKategori != "" {
		query = query.Where("sub_kategori ILIKE ?", "%"+subKategori+"%")
	}
	if tahap != "" {
		query = query.Where("tahap ILIKE ?", "%"+tahap+"%")
	}
	if search != "" {
		query = query.Where("tim1 ILIKE ? OR tim2 ILIKE ? OR hasil ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghitung total data",
		})
		return
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data pertandingan knockout",
		})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pertandingan knockout berhasil diambil",
		"data":    matches,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	})
}

func GetKnockoutMatchByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	var match models.KnockoutMatch
	if err := config.DB.First(&match, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Pertandingan knockout tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data pertandingan knockout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pertandingan knockout berhasil diambil",
		"data":    match,
	})
}

func UpdateKnockoutMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	var req UpdateKnockoutMatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	var match models.KnockoutMatch
	if err := config.DB.First(&match, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Pertandingan knockout tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data pertandingan knockout",
		})
		return
	}

	if req.Kategori != "" {
		match.Kategori = strings.TrimSpace(req.Kategori)
	}
	if req.SubKategori != "" {
		match.SubKategori = strings.TrimSpace(req.SubKategori)
	}
	if req.Tim1 != "" {
		match.Tim1 = strings.TrimSpace(req.Tim1)
	}
	if req.Tim2 != "" {
		match.Tim2 = strings.TrimSpace(req.Tim2)
	}
	if req.Hasil != "" {
		match.Hasil = strings.TrimSpace(req.Hasil)
	}
	if req.Tahap != "" {
		match.Tahap = strings.TrimSpace(req.Tahap)
	}

	if err := config.DB.Save(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate data pertandingan knockout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pertandingan knockout berhasil diupdate",
		"data":    match,
	})
}

func DeleteKnockoutMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	var match models.KnockoutMatch
	if err := config.DB.First(&match, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Pertandingan knockout tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data pertandingan knockout",
		})
		return
	}

	if err := config.DB.Delete(&match).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus pertandingan knockout",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Pertandingan knockout berhasil dihapus",
	})
}

func GetKnockoutStats(c *gin.Context) {
	var stats struct {
		Total      int64            `json:"total"`
		ByKategori map[string]int64 `json:"by_kategori"`
		ByTahap    map[string]int64 `json:"by_tahap"`
	}

	config.DB.Model(&models.KnockoutMatch{}).Count(&stats.Total)

	var kategoriResults []struct {
		Kategori string
		Count    int64
	}
	config.DB.Model(&models.KnockoutMatch{}).
		Select("kategori, COUNT(*) as count").
		Group("kategori").
		Scan(&kategoriResults)

	stats.ByKategori = make(map[string]int64)
	for _, result := range kategoriResults {
		stats.ByKategori[result.Kategori] = result.Count
	}

	var tahapResults []struct {
		Tahap string
		Count int64
	}
	config.DB.Model(&models.KnockoutMatch{}).
		Select("tahap, COUNT(*) as count").
		Group("tahap").
		Scan(&tahapResults)

	stats.ByTahap = make(map[string]int64)
	for _, result := range tahapResults {
		stats.ByTahap[result.Tahap] = result.Count
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Statistik pertandingan knockout berhasil diambil",
		"data":    stats,
	})
}

func GetKnockoutByKategori(c *gin.Context) {
	kategori := c.Param("kategori")
	if kategori == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Kategori tidak boleh kosong",
		})
		return
	}

	var matches []models.KnockoutMatch
	if err := config.DB.Where("kategori ILIKE ?", "%"+kategori+"%").Order("created_at DESC").Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data pertandingan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pertandingan berhasil diambil",
		"data":    matches,
	})
}

func GetKnockoutByTahap(c *gin.Context) {
	tahap := c.Param("tahap")
	if tahap == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Tahap tidak boleh kosong",
		})
		return
	}

	var matches []models.KnockoutMatch
	if err := config.DB.Where("tahap ILIKE ?", "%"+tahap+"%").Order("created_at DESC").Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data pertandingan",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data pertandingan berhasil diambil",
		"data":    matches,
	})
}

func GetKnockoutStanding(c *gin.Context) {
	kategori := c.Query("kategori")
	subKategori := c.Query("sub_kategori")

	query := config.DB.Model(&models.KnockoutMatch{})

	if kategori != "" {
		query = query.Where("kategori ILIKE ?", "%"+kategori+"%")
	}
	if subKategori != "" {
		query = query.Where("sub_kategori ILIKE ?", "%"+subKategori+"%")
	}

	var matches []models.KnockoutMatch
	if err := query.Order("tahap ASC, created_at DESC").Find(&matches).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data standing knockout",
		})
		return
	}

	standingByTahap := make(map[string][]models.KnockoutMatch)
	for _, match := range matches {
		standingByTahap[match.Tahap] = append(standingByTahap[match.Tahap], match)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data standing knockout berhasil diambil",
		"data":    standingByTahap,
	})
}
