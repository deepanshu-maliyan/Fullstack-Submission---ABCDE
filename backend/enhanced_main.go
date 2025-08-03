package main

import (
	"ecommerce-backend/database"
	"ecommerce-backend/handlers"
	"ecommerce-backend/middleware"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize in-memory database
	database.Connect()

	// Set Gin mode to release for production-like behavior
	// gin.SetMode(gin.ReleaseMode) // Uncomment for production

	// Initialize Gin router
	r := gin.Default()

	// Enhanced CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://127.0.0.1:3000",
			"http://192.168.29.248:3000",
			"https://your-domain.com", // Add your production domain
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization", "Accept",
			"X-Requested-With", "Access-Control-Request-Method",
			"Access-Control-Request-Headers",
		},
		ExposeHeaders: []string{
			"Content-Length", "X-Total-Count", "X-Total-Pages",
			"X-Current-Page", "X-Per-Page",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Request logging middleware
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	// Recovery middleware with custom handler
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Internal server error",
				"meta": gin.H{
					"timestamp": time.Now().Format(time.RFC3339),
					"version":   "v1.0",
				},
			})
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	// Rate limiting middleware (simple in-memory implementation)
	r.Use(middleware.RateLimit(100, time.Minute)) // 100 requests per minute

	// Security headers middleware
	r.Use(middleware.SecurityHeaders())

	// Request ID middleware
	r.Use(middleware.RequestID())

	// Serve static files (assets/images) with cache headers
	r.Static("/assets", "../assets")
	r.StaticFile("/favicon.ico", "../assets/favicon.ico")

	// API version prefix
	api := r.Group("/api/v1")
	{
		// Public endpoints
		public := api.Group("/")
		{
			public.POST("/login", handlers.LoginUser)
			public.GET("/items", handlers.GetItems)
			public.GET("/health", handlers.HealthCheck)
		}

		// Protected endpoints
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// User management
			protected.GET("/users", handlers.GetUsers)
			protected.GET("/profile", middleware.GetUserProfile())

			// Item management
			protected.POST("/items", handlers.CreateItem)

			// Cart management
			protected.POST("/carts", handlers.AddToCart)
			protected.GET("/carts", handlers.GetCarts)
			protected.GET("/carts/user", handlers.GetUserCart)
			protected.GET("/carts/:id", handlers.GetCartByID)
			protected.DELETE("/carts/clear", middleware.ClearCart())

			// Order management
			protected.POST("/orders", handlers.CreateOrder)
			protected.GET("/orders", handlers.GetOrders)
			protected.GET("/orders/user", handlers.GetUserOrders)
		}
	}

	// Legacy endpoints (for backward compatibility)
	r.POST("/users/login", handlers.LoginUser)
	r.GET("/items", handlers.GetItems)

	// Protected legacy routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	{
		// User routes
		auth.GET("/users", handlers.GetUsers)

		// Item routes
		auth.POST("/items", handlers.CreateItem)

		// Cart routes
		auth.POST("/carts", handlers.AddToCart)
		auth.GET("/carts", handlers.GetCarts)
		auth.GET("/carts/user", handlers.GetUserCart)
		auth.GET("/carts/:id", handlers.GetCartByID)

		// Order routes
		auth.POST("/orders", handlers.CreateOrder)
		auth.GET("/orders", handlers.GetOrders)
		auth.GET("/orders/user", handlers.GetUserOrders)
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "Endpoint not found",
			"meta": gin.H{
				"timestamp": time.Now().Format(time.RFC3339),
				"version":   "v1.0",
			},
		})
	})

	// Graceful server startup
	log.Println("🚀 ModernStore Backend Server Starting...")
	log.Println("📍 Server URL: http://localhost:8080")
	log.Println("📊 Health Check: http://localhost:8080/api/v1/health")
	log.Println("📁 Static Assets: http://localhost:8080/assets")
	log.Println("🔒 Admin Login: POST /api/v1/login")
	log.Println("📦 Items API: GET /api/v1/items")
	log.Println("✨ Enhanced features: Rate limiting, Security headers, Request logging")
	log.Println("🎯 Ready to handle requests!")

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
