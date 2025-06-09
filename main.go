package main

import (
	"os"
	"time"

	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/models"
	"github.com/ImamIryunullah/BE-PEP/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	corsConfig := cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	router.Use(cors.New(corsConfig))

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	config.ConnectDatabase()

	config.DB.AutoMigrate(
		&models.DaftarUser{},
		&models.Berita{},
		&models.ParticipantRegistration{},
		&models.Funrun{}, 
	)

	if _, err := os.Stat("uploads"); os.IsNotExist(err) {
		os.Mkdir("uploads", 0755)
	}

	router.Static("/uploads", "./uploads")

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
			"time":    time.Now(),
		})
	})

	routes.SetupRoutes(router)

	router.Run("0.0.0.0:8080")
}
