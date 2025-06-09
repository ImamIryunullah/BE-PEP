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
		beritaGroup.GET("/", controllers.GetAllBerita)
		beritaGroup.DELETE("/:id", controllers.DeleteBerita)
	}

	router.POST("/daftar", controllers.SubmitParticipantRegistration)
	router.GET("/daftar", controllers.GetAllPeserta)
	router.PUT("/daftar/:id", controllers.EditParticipantRegistration)
	router.PUT("/daftar/:id/status", controllers.UpdateParticipantStatus)
	router.GET("/daftar/:id", controllers.GetParticipantById)

	funrunGroup := router.Group("/funrun")
	{

		funrunGroup.POST("/peserta", controllers.CreatePesertaFunrun)           // Create peserta funrun
		funrunGroup.GET("/peserta", controllers.GetAllPesertaFunrun)            // Get all peserta with pagination & filters
		funrunGroup.GET("/peserta/:id", controllers.GetPesertaByID)             // Get peserta by ID
		funrunGroup.PUT("/peserta/:id", controllers.UpdatePesertaFunrun)        // Update peserta
		funrunGroup.PUT("/peserta/:id/status", controllers.UpdatePesertaStatus) // Update status only
		funrunGroup.DELETE("/peserta/:id", controllers.DeletePesertaFunrun)     // Delete peserta (soft delete)

		funrunGroup.GET("/stats", controllers.GetPesertaStats)                      // Get peserta statistics
		funrunGroup.GET("/kontingen/:kontingen", controllers.GetPesertaByKontingen) // Get peserta by kontingen
	}
}
