package handlers

import (
	"ecommerce-backend/database"
	"ecommerce-backend/models"
	"ecommerce-backend/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Enhanced response structure
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Timestamp string `json:"timestamp"`
	RequestID string `json:"request_id,omitempty"`
	Version   string `json:"version"`
}

// Enhanced login with better validation and logging
func LoginUser(c *gin.Context) {
	var loginRequest struct {
		Username string `json:"username" binding:"required,min=3,max=50"`
		Password string `json:"password" binding:"required,min=6"`
	}

	// Enhanced request validation
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		log.Printf("Login validation error: %v", err)
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid input: " + err.Error(),
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Sanitize input
	loginRequest.Username = strings.TrimSpace(strings.ToLower(loginRequest.Username))
	
	log.Printf("Login attempt for user: %s from IP: %s", loginRequest.Username, c.ClientIP())

	// Find user
	var user models.User
	if err := database.DB.Where("username = ?", loginRequest.Username).First(&user).Error; err != nil {
		log.Printf("User not found: %s", loginRequest.Username)
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Invalid credentials",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		log.Printf("Invalid password for user: %s", loginRequest.Username)
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "Invalid credentials",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		log.Printf("Token generation error for user %s: %v", loginRequest.Username, err)
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to generate token",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	log.Printf("Successful login for user: %s", loginRequest.Username)

	// Enhanced response
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("Welcome back, %s!", user.Username),
		Data: gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"role":     "admin",
			},
			"expires_in": "24h",
		},
		Meta: &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}

// Enhanced GetItems with pagination and filtering
func GetItems(c *gin.Context) {
	// Parse query parameters
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "20")
	category := c.Query("category")
	search := c.Query("search")
	status := c.Query("status")

	page, _ := strconv.Atoi(pageStr)
	limit, _ := strconv.Atoi(limitStr)
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	// Build query
	query := database.DB.Model(&models.Item{})

	if category != "" && category != "all" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(category)+"%")
	}

	if search != "" {
		query = query.Where("LOWER(name) LIKE ? OR LOWER(status) LIKE ?", 
			"%"+strings.ToLower(search)+"%", "%"+strings.ToLower(search)+"%")
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get items with pagination
	var items []models.Item
	if err := query.Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		log.Printf("Error fetching items: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to fetch items",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Calculate pagination info
	totalPages := (total + int64(limit) - 1) / int64(limit)

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("Found %d items", len(items)),
		Data:    items,
		Meta: &Meta{
			Timestamp: time.Now().Format(time.RFC3339),
			Version:   "v1.0",
		},
	})

	// Add pagination info to header
	c.Header("X-Total-Count", strconv.FormatInt(total, 10))
	c.Header("X-Total-Pages", strconv.FormatInt(totalPages, 10))
	c.Header("X-Current-Page", strconv.Itoa(page))
	c.Header("X-Per-Page", strconv.Itoa(limit))
}

// Enhanced CreateItem with validation
func CreateItem(c *gin.Context) {
	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		log.Printf("Item creation validation error: %v", err)
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid item data: " + err.Error(),
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Validate required fields
	if strings.TrimSpace(item.Name) == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Item name is required",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Set default status if not provided
	if item.Status == "" {
		item.Status = "available"
	}

	// Create item
	if err := database.DB.Create(&item).Error; err != nil {
		log.Printf("Error creating item: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to create item",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	log.Printf("Item created successfully: %s (ID: %d)", item.Name, item.ID)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: "Item created successfully",
		Data:    item,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}

// Enhanced AddToCart with duplicate checking
func AddToCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "User not authenticated",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	var request struct {
		ItemID uint `json:"item_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request: " + err.Error(),
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Check if item exists
	var item models.Item
	if err := database.DB.First(&item, request.ItemID).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Item not found",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Check if item is available
	if item.Status != "available" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Item is not available",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Find or create cart
	var cart models.Cart
	if err := database.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		cart = models.Cart{UserID: userID.(uint)}
		database.DB.Create(&cart)
	}

	// Check if item already in cart
	var existingCartItem models.CartItem
	if err := database.DB.Where("cart_id = ? AND item_id = ?", cart.ID, request.ItemID).First(&existingCartItem).Error; err == nil {
		c.JSON(http.StatusConflict, Response{
			Success: false,
			Error:   "Item already in cart",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Add item to cart
	cartItem := models.CartItem{
		CartID: cart.ID,
		ItemID: request.ItemID,
	}

	if err := database.DB.Create(&cartItem).Error; err != nil {
		log.Printf("Error adding item to cart: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to add item to cart",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	log.Printf("Item %d added to cart for user %v", request.ItemID, userID)

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: fmt.Sprintf("'%s' added to cart successfully", item.Name),
		Data:    cartItem,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}

// Enhanced GetUserCart with item details
func GetUserCart(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "User not authenticated",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	var cart models.Cart
	if err := database.DB.Where("user_id = ?", userID).Preload("CartItems.Item").First(&cart).Error; err != nil {
		// Return empty cart if not found
		c.JSON(http.StatusOK, Response{
			Success: true,
			Message: "Cart is empty",
			Data: gin.H{
				"id":         0,
				"user_id":    userID,
				"cart_items": []models.CartItem{},
				"total_items": 0,
			},
			Meta: &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("Cart contains %d items", len(cart.CartItems)),
		Data: gin.H{
			"id":          cart.ID,
			"user_id":     cart.UserID,
			"cart_items":  cart.CartItems,
			"total_items": len(cart.CartItems),
		},
		Meta: &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}

// Enhanced CreateOrder with better validation
func CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "User not authenticated",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Find user's cart
	var cart models.Cart
	if err := database.DB.Where("user_id = ?", userID).Preload("CartItems.Item").First(&cart).Error; err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "No cart found",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	if len(cart.CartItems) == 0 {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Cart is empty",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Create order
	order := models.Order{
		UserID: userID.(uint),
		Status: "pending",
	}

	if err := database.DB.Create(&order).Error; err != nil {
		log.Printf("Error creating order: %v", err)
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to create order",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	// Clear cart after successful order
	database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	
	log.Printf("Order %d created successfully for user %v with %d items", order.ID, userID, len(cart.CartItems))

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Message: fmt.Sprintf("Order #%d created successfully with %d items", order.ID, len(cart.CartItems)),
		Data: gin.H{
			"order_id":    order.ID,
			"status":      order.Status,
			"items_count": len(cart.CartItems),
			"created_at":  order.CreatedAt,
		},
		Meta: &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}

// Enhanced GetUserOrders with pagination
func GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "User not authenticated",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	var orders []models.Order
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		log.Printf("Error fetching orders for user %v: %v", userID, err)
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to fetch orders",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("Found %d orders", len(orders)),
		Data:    orders,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}

// Health check endpoint
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Server is healthy",
		Data: gin.H{
			"status":    "ok",
			"timestamp": time.Now().Format(time.RFC3339),
			"version":   "v1.0",
			"uptime":    "24/7",
		},
		Meta: &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}

// Get all users (admin only)
func GetUsers(c *gin.Context) {
	var users []models.User
	if err := database.DB.Select("id, username, created_at, updated_at").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to fetch users",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("Found %d users", len(users)),
		Data:    users,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}

// Get all carts (admin only)
func GetCarts(c *gin.Context) {
	var carts []models.Cart
	if err := database.DB.Preload("CartItems.Item").Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to fetch carts",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("Found %d carts", len(carts)),
		Data:    carts,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}

// Get cart by ID (admin only)
func GetCartByID(c *gin.Context) {
	id := c.Param("id")
	
	var cart models.Cart
	if err := database.DB.Where("id = ?", id).Preload("CartItems.Item").First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Cart not found",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: "Cart found",
		Data:    cart,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}

// Get all orders (admin only)
func GetOrders(c *gin.Context) {
	var orders []models.Order
	if err := database.DB.Order("created_at DESC").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to fetch orders",
			Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Message: fmt.Sprintf("Found %d orders", len(orders)),
		Data:    orders,
		Meta:    &Meta{Timestamp: time.Now().Format(time.RFC3339), Version: "v1.0"},
	})
}
