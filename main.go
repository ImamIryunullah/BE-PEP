package main

import (
	"log"

	"github.com/ImamIryunullah/BE-PEP/config"
	"github.com/ImamIryunullah/BE-PEP/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CheckDatabaseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.DB == nil {
			c.JSON(500, gin.H{"error": "Database is not connected"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func setupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://192.168.1.53:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "csrf-token", "CSRF-TOKEN", "X-Requested-With"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(gin.Recovery())
	routes.SetupRoutes(router)
	return router
}
func main() {
	gin.SetMode(gin.ReleaseMode)
	config.ConnectDatabase()
	r := setupRouter()
	r.Use(CheckDatabaseMiddleware())
	if err := r.Run("192.168.1.53:3000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}

}
