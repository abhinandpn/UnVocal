package main

import (
	"log"

	"github.com/abhinandpn/UnVocal/services/user-service/config"
	"github.com/abhinandpn/UnVocal/services/user-service/db"
	"github.com/abhinandpn/UnVocal/services/user-service/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Connect to PostgreSQL
	if err := db.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal(err)
	}

	// Create Gin router
	r := gin.Default()

	// Register routes
	routes.SetupRoutes(r, db.DB)

	// Health check
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "User Service Running 🚀",
		})
	})

	log.Printf("🚀 Server is running on http://localhost:%s\n", cfg.Port)

	// Start server
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
