package main

import (
	"log"

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

	port := cfg.Port
	log.Printf("Starting server on port %s", port)

	addr := ":" + port
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
