package main

import (
	"ecommerce-backend/database"
	"ecommerce-backend/handlers"
	"ecommerce-backend/middleware"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize in-memory database
	database.Connect()

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://192.168.29.248:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve static files (assets/images)
	r.Static("/assets", "../assets")

	// Public routes
	r.POST("/users/login", handlers.LoginUser)
	r.GET("/items", handlers.GetItems)

	// Protected routes
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

	log.Println("Server starting on http://localhost:8080")
	r.Run(":8080")
}
