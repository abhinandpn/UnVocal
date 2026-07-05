package main

import (
	"log"

	_ "github.com/abhinandpn/UnVocal/services/user-service/docs"

	"github.com/abhinandpn/UnVocal/services/user-service/config"
	"github.com/abhinandpn/UnVocal/services/user-service/db"
	"github.com/abhinandpn/UnVocal/services/user-service/routes"
	"github.com/gin-gonic/gin"
)

// @title User Service API - UnVocal
// @version 1.0
// @description This is the API documentation for the User Service of the UnVocal application.
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to PostgreSQL
	if err := db.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal(err)
	}

	// Create Gin router
	router := gin.Default()

	// Register Swagger
	routes.SetupSwagger(router)

	// Register application routes
	routes.SetupRoutes(router, db.DB)

	// Health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "User Service Running 🚀",
		})
	})

	log.Printf("🚀 Server is running on http://localhost:%s", cfg.Port)
	log.Printf("📚 Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)

	// Start server
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
