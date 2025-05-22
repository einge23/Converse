package main

import (
	"log"
	"os"

	"converse/config"

	"github.com/gin-gonic/gin"
)

func main() {

	cfg := config.New()

	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalf("PORT is not set")
		return
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
