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

	router.GET("/registrations/user/:user_id", controllers.GetRegistrationsByUserID)
	router.GET("/users/:user_id/registrations", controllers.GetUserWithRegistrations)

	beritaGroup := router.Group("/berita")
	{
		beritaGroup.POST("/", controllers.CreateBerita)
		beritaGroup.PUT("/:id", controllers.UpdateBerita)
		router.GET("/berita", controllers.GetAllBerita)
	}

	router.POST("/daftar", controllers.SubmitParticipantRegistration)
	router.GET("/daftar", controllers.GetAllPeserta)
	router.PUT("/daftar/:id", controllers.EditParticipantRegistration)
	router.PUT("/daftar/:id/status", controllers.UpdateParticipantStatus)
	router.GET("/daftar/:id", controllers.GetParticipantById)

}
