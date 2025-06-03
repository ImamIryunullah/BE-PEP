package routes

import (
	"github.com/ImamIryunullah/BE-PEP/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.POST("/register", controllers.RegisterPeserta)
	router.POST("/login", controllers.Login)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	beritaGroup := router.Group("/berita")
	{
		beritaGroup.POST("/", controllers.CreateBerita)
		beritaGroup.PUT("/:id", controllers.UpdateBerita)
		router.GET("/berita", controllers.GetAllBerita)
		beritaGroup.DELETE("/:id", controllers.DeleteBerita)
	}
}
