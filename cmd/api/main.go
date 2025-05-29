package main

import (
	"log"

	"converse/internal/config"
	"converse/internal/db"
	"converse/internal/handlers"
	"converse/internal/middleware"
	"converse/migrations"

	"github.com/gin-contrib/cors"
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

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "https://converse-ui-development.up.railway.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "x-device-id", "x-session-id"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

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
        authHandler := handlers.NewAuthHandler()
		friendRequestHandler := handlers.NewFriendRequestHandler()
		friendshipHandler := handlers.NewFriendshipHandler()

        // Auth routes
        auth := v1.Group("/auth")
        {
            auth.POST("/register", authHandler.Register)
            auth.POST("/login", authHandler.Login)
            auth.GET("/validate-session", authHandler.ValidateSession)
        }

        // Protected routes
        protected := v1.Group("/")
        protected.Use(middleware.AuthMiddleware())
        {
            // Auth management routes
            auth := protected.Group("/auth")
            {
                auth.POST("/logout", authHandler.Logout)
            }
			
			friend_requests := protected.Group("/friend-requests")
			{
				friend_requests.POST("/", friendRequestHandler.CreateFriendRequest)
			}
			
			friend_requests.Use(middleware.AuthMiddleware())
			{
				friend_requests.GET("/", friendRequestHandler.GetUserFriendRequests)
				friend_requests.PUT("/:friend_request_id/accept", friendRequestHandler.AcceptFriendRequest)
				friend_requests.PUT("/:friend_request_id/decline", friendRequestHandler.DeclineFriendRequest)
			}

			friend_requests.Use(middleware.OwnResourceMiddleware())
			{
				friend_requests.POST("/:friend_request_id/accept", friendRequestHandler.AcceptFriendRequest)
				friend_requests.POST("/:friend_request_id/decline", friendRequestHandler.DeclineFriendRequest)
			}

            friendships := protected.Group("/friends").Use(middleware.AuthMiddleware())
			{
				friendships.GET("/", friendshipHandler.GetFriends)
			}
        }
    }
}