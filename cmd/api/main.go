package main

import (
	"log"

	"converse/internal/config"
	"converse/internal/db"
	"converse/internal/handlers"
	"converse/internal/middleware"
	"converse/migrations"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.New()

	if err := db.Init(cfg); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }
    defer db.Close()

	if err := migrations.RunMigrations(); err != nil {
        log.Fatalf("Failed to run migrations: %v", err)
    }

	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	setupRoutes(r)

	port := cfg.Port
	log.Printf("Starting server on port %s", port)

	addr := ":" + port
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}


}


func setupRoutes(r *gin.Engine) {
    // API v1 routes
    v1 := r.Group("/api/v1")
    {
        // Auth routes
        auth := v1.Group("/auth")
        {
            authHandler := handlers.NewAuthHandler()
            auth.POST("/register", authHandler.Register)
        }

        // Protected routes
        protected := v1.Group("/")
        protected.Use(middleware.AuthMiddleware())
        {
            // Add protected routes here
        }
    }
}