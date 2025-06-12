package routes

import (
	"github.com/ImamIryunullah/BE-PEP/controllers"
	"github.com/ImamIryunullah/BE-PEP/handlers"
	"github.com/ImamIryunullah/BE-PEP/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	statis := router.Group("/")
	berita := router.Group("/")
	protected := router.Group("/api")
	protected.Use(middleware.VerifyJWT())
	router.POST("/api/register", controllers.RegisterPeserta)
	router.GET("/api/register", controllers.GetAllAkun)
	router.POST("/api/login", handlers.Login)
	router.POST("/api/logout", handlers.Logout)
	protected.GET("/datauser", handlers.GetUserDataAuth)
	router.GET("/api/registrations/user/:user_id", controllers.GetRegistrationsByUserID)
	router.GET("/api/users/:user_id/registrations", controllers.GetUserWithRegistrations)
	beritaGroup := router.Group("/api/berita")
	{
		beritaGroup.POST("/", controllers.CreateBerita)
		beritaGroup.GET("/", controllers.GetAllBerita)
		beritaGroup.GET("/:id", controllers.GetBeritaByID)
		beritaGroup.PUT("/:id", controllers.UpdateBerita)
		beritaGroup.DELETE("/:id", controllers.DeleteBerita)
	}
	protected.POST("/daftar", controllers.SubmitParticipantRegistration)
	protected.GET("/daftar-peserta", controllers.GetAllPeserta)        //admin
	router.GET("api/daftar-list", controllers.GetAllPesertaList)       //admin
	protected.GET("/daftar-akun", controllers.GetParticipantsByUserID) //user
	protected.PUT("/daftar/:id", controllers.EditParticipantRegistration)
	protected.PUT("/daftar/:id/status", controllers.UpdateParticipantStatus)
	// protected.GET("/daftar/:id", controllers.GetParticipantById)

	funrunGroup := router.Group("/api/funrun")
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
	knockoutGroup := router.Group("/api/knockout")
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
	assets := router.Group("/")
	assets.Use(middleware.CacheControlMiddleware())
	assets.Use(gzip.Gzip(gzip.BestSpeed))
	assets.Static("/assets", "./assets")

	statis.Use(middleware.VerifyJWT())
	statis.Static("/uploads", "./uploads")

	berita.Static("/berita", "./berita")
}
