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

	// 1. CONNECT DB FIRST
	if err := db.Connect(cfg.DatabaseURL); err != nil {
		log.Fatal(err)
	}
	// 2. RUN MIGRATIONS
	err := db.RunMigrations(db.DB)
	if err != nil {
		log.Fatal(err)
	}

	// 3. CREATE ROUTER
	r := gin.Default()

	// 4. REGISTER ROUTES
	routes.SetupRoutes(r, db.DB)

	// 5. HEALTH CHECK
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "User Service Running 🚀",
		})
	})

	log.Printf("🚀 Server is running on http://localhost:%s\n", cfg.Port)

	// 6. START SERVER
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
