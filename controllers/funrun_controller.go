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

// Request structs
type CreatePesertaFunrunRequest struct {
	Nama      string `json:"nama" binding:"required,max=100"`
	Email     string `json:"email" binding:"required,email,max=100"`
	Kontingen string `json:"kontingen" binding:"required,max=100"`
	Size      string `json:"size" binding:"required,max=10"`
}

type UpdatePesertaFunrunRequest struct {
	Nama      string `json:"nama" binding:"max=100"`
	Email     string `json:"email" binding:"email,max=100"`
	Kontingen string `json:"kontingen" binding:"max=100"`
	Size      string `json:"size" binding:"max=10"`
	Status    string `json:"status" binding:"oneof=pending approved rejected"`
}

type UpdateStatusFunRunRequest struct {
	Status string `json:"status" binding:"required,oneof=pending approved rejected"`
}

func CreatePesertaFunrun(c *gin.Context) {
	var req CreatePesertaFunrunRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Nama = strings.TrimSpace(req.Nama)
	req.Kontingen = strings.TrimSpace(req.Kontingen)
	req.Size = strings.TrimSpace(req.Size)

	existingPeserta, err := GetPesertaByEmailFromModel(req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal memeriksa data peserta",
		})
		return
	}

	if existingPeserta != nil {
		c.JSON(http.StatusConflict, gin.H{
			"success": false,
			"message": "Email sudah terdaftar",
		})
		return
	}

	peserta := models.Funrun{
		Nama:      req.Nama,
		Email:     req.Email,
		Kontingen: req.Kontingen,
		Size:      req.Size,
		Status:    "pending",
	}

	if err := config.DB.Create(&peserta).Error; err != nil {
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Email sudah terdaftar",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menyimpan data peserta",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Peserta berhasil ditambahkan",
		"data":    peserta,
	})
}

func GetAllPesertaFunrun(c *gin.Context) {
	var peserta []models.Funrun
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	kontingen := c.Query("kontingen")
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := config.DB.Model(&models.Funrun{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if kontingen != "" {
		query = query.Where("kontingen ILIKE ?", "%"+kontingen+"%")
	}
	if search != "" {
		query = query.Where("nama ILIKE ? OR email ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghitung total data",
		})
		return
	}

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&peserta).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data peserta",
		})
		return
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data peserta berhasil diambil",
		"data":    peserta,
		"pagination": gin.H{
			"current_page": page,
			"total_pages":  totalPages,
			"total_items":  total,
			"per_page":     limit,
		},
	})
}

func GetPesertaByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	var peserta models.Funrun
	if err := config.DB.First(&peserta, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Peserta tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data peserta",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data peserta berhasil diambil",
		"data":    peserta,
	})
}

func UpdatePesertaFunrun(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	var req UpdatePesertaFunrunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	var peserta models.Funrun
	if err := config.DB.First(&peserta, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Peserta tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data peserta",
		})
		return
	}

	if req.Nama != "" {
		peserta.Nama = strings.TrimSpace(req.Nama)
	}
	if req.Email != "" {
		newEmail := strings.ToLower(strings.TrimSpace(req.Email))
		if newEmail != peserta.Email {
			existingPeserta, err := GetPesertaByEmailFromModel(newEmail)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Gagal memeriksa email",
				})
				return
			}
			if existingPeserta != nil {
				c.JSON(http.StatusConflict, gin.H{
					"success": false,
					"message": "Email sudah digunakan peserta lain",
				})
				return
			}
		}
		peserta.Email = newEmail
	}
	if req.Kontingen != "" {
		peserta.Kontingen = strings.TrimSpace(req.Kontingen)
	}
	if req.Size != "" {
		peserta.Size = strings.TrimSpace(req.Size)
	}
	if req.Status != "" {
		peserta.Status = req.Status
	}

	if err := config.DB.Save(&peserta).Error; err != nil {
		if isDuplicateError(err) {
			c.JSON(http.StatusConflict, gin.H{
				"success": false,
				"message": "Email sudah digunakan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate data peserta",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data peserta berhasil diupdate",
		"data":    peserta,
	})
}

func UpdatePesertaStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Data tidak valid",
			"errors":  err.Error(),
		})
		return
	}

	var peserta models.Funrun
	if err := config.DB.First(&peserta, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Peserta tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data peserta",
		})
		return
	}

	peserta.Status = req.Status

	if err := config.DB.Save(&peserta).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengupdate status peserta",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Status peserta berhasil diupdate",
		"data":    peserta,
	})
}

// DeletePeserta - Menghapus peserta (soft delete)
func DeletePesertaFunrun(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID tidak valid",
		})
		return
	}

	var peserta models.Funrun
	if err := config.DB.First(&peserta, uint(id)).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "Peserta tidak ditemukan",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data peserta",
		})
		return
	}

	if err := config.DB.Delete(&peserta).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal menghapus peserta",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Peserta berhasil dihapus",
	})
}

// GetPesertaStats - Mengambil statistik peserta
func GetPesertaStats(c *gin.Context) {
	var stats struct {
		Total    int64 `json:"total"`
		Pending  int64 `json:"pending"`
		Approved int64 `json:"approved"`
		Rejected int64 `json:"rejected"`
	}

	// Total peserta
	config.DB.Model(&models.Funrun{}).Count(&stats.Total)

	// Pending
	config.DB.Model(&models.Funrun{}).Where("status = ?", "pending").Count(&stats.Pending)

	// Approved
	config.DB.Model(&models.Funrun{}).Where("status = ?", "approved").Count(&stats.Approved)

	// Rejected
	config.DB.Model(&models.Funrun{}).Where("status = ?", "rejected").Count(&stats.Rejected)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Statistik peserta berhasil diambil",
		"data":    stats,
	})
}

// Helper functions
func GetPesertaByEmailFromModel(email string) (*models.Funrun, error) {
	var peserta models.Funrun
	err := config.DB.Where("email = ?", strings.ToLower(strings.TrimSpace(email))).First(&peserta).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &peserta, nil
}

// GetPesertaByKontingen - Mengambil peserta berdasarkan kontingen
func GetPesertaByKontingen(c *gin.Context) {
	kontingen := c.Param("kontingen")
	if kontingen == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Nama kontingen tidak boleh kosong",
		})
		return
	}

	var peserta []models.Funrun
	if err := config.DB.Where("kontingen ILIKE ?", "%"+kontingen+"%").Order("created_at DESC").Find(&peserta).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Gagal mengambil data peserta",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Data peserta berhasil diambil",
		"data":    peserta,
	})
}
