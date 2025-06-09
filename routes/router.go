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
		beritaGroup.GET("/", controllers.GetAllBerita)
		beritaGroup.GET("/:id", controllers.GetBeritaByID)
		beritaGroup.PUT("/:id", controllers.UpdateBerita)
		beritaGroup.DELETE("/:id", controllers.DeleteBerita)
	}
	router.POST("/daftar", controllers.SubmitParticipantRegistration)
	router.GET("/daftar", controllers.GetAllPeserta)
	router.PUT("/daftar/:id", controllers.EditParticipantRegistration)
	router.PUT("/daftar/:id/status", controllers.UpdateParticipantStatus)
	router.GET("/daftar/:id", controllers.GetParticipantById)

	funrunGroup := router.Group("/funrun")
	{
		funrunGroup.POST("/peserta", controllers.CreatePesertaFunrun)
		funrunGroup.GET("/peserta", controllers.GetAllPesertaFunrun)
		funrunGroup.GET("/peserta/:id", controllers.GetPesertaByID)
		funrunGroup.PUT("/peserta/:id", controllers.UpdatePesertaFunrun)
		funrunGroup.PUT("/peserta/:id/status", controllers.UpdatePesertaStatus)
		funrunGroup.DELETE("/peserta/:id", controllers.DeletePesertaFunrun)

		funrunGroup.GET("/stats", controllers.GetPesertaStats)
		funrunGroup.GET("/kontingen/:kontingen", controllers.GetPesertaByKontingen)
	}
	knockoutGroup := router.Group("/knockout")
	{
		knockoutGroup.POST("/", controllers.CreateKnockoutMatch)
		knockoutGroup.GET("/", controllers.GetAllKnockoutMatches)
		knockoutGroup.GET("/:id", controllers.GetKnockoutMatchByID)
		knockoutGroup.PUT("/:id", controllers.UpdateKnockoutMatch)
		knockoutGroup.DELETE("/:id", controllers.DeleteKnockoutMatch)

		knockoutGroup.GET("/stats", controllers.GetKnockoutStats)
		knockoutGroup.GET("/kategori/:kategori", controllers.GetKnockoutByKategori)
		knockoutGroup.GET("/tahap/:tahap", controllers.GetKnockoutByTahap)
		knockoutGroup.GET("/standing", controllers.GetKnockoutStanding)
	}
}
