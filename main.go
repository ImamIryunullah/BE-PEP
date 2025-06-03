package main

import (
	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/models"
	"github.com/ImamIryunullah/BE-PEP/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.Use(cors.Default())
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Peserta{}, &models.Berita{})
	routes.SetupRoutes(router)
	router.Run("0.0.0.0:8080")

}
