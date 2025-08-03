package handlers

import (
	"ecommerce-backend/database"
	"ecommerce-backend/models"
	"ecommerce-backend/utils"
	"fmt"
	"net/http"
	"time"
	"golang.org/x/crypto/bcrypt"
	"github.com/gin-gonic/gin"
)

type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

func CreateUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Mutex.Lock()
	defer database.DB.Mutex.Unlock()

	// Check if user already exists
	for _, user := range database.DB.Users {
		if user.Username == req.Username {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	userID := database.DB.GetNextID()
	user := &models.User{
		ID:        userID,
		Username:  req.Username,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Create a cart for the user
	cartID := database.DB.GetNextID()
	cart := &models.Cart{
		ID:        cartID,
		UserID:    userID,
		Name:      "Default Cart",
		Status:    "active",
		CreatedAt: time.Now(),
		CartItems: []models.CartItem{},
	}

	// Store in database
	database.DB.Carts[cartID] = cart
	user.CartID = cartID
	database.DB.Users[userID] = user

	// Remove password from response
	responseUser := *user
	responseUser.Password = ""
	c.JSON(http.StatusCreated, responseUser)
}

func LoginUser(c *gin.Context) {
	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Mutex.RLock()
	defer database.DB.Mutex.RUnlock()

	var user *models.User
	for _, u := range database.DB.Users {
		if u.Username == req.Username {
			user = u
			break
		}
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	// Generate token
	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Update user token
	user.Token = token

	// Remove password from response
	responseUser := *user
	responseUser.Password = ""

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User:  responseUser,
	})
}

func GetUsers(c *gin.Context) {
	database.DB.Mutex.RLock()
	defer database.DB.Mutex.RUnlock()

	var users []models.User
	for _, user := range database.DB.Users {
		responseUser := *user
		responseUser.Password = "" // Remove password
		users = append(users, responseUser)
	}

	c.JSON(http.StatusOK, users)
}

func GetItems(c *gin.Context) {
	database.DB.Mutex.RLock()
	defer database.DB.Mutex.RUnlock()

	var items []models.Item
	for _, item := range database.DB.Items {
		if item.Status == "active" {
			items = append(items, *item)
		}
	}

	c.JSON(http.StatusOK, items)
}

type ItemRequest struct {
	Name   string `json:"name" binding:"required"`
	Status string `json:"status"`
}

func CreateItem(c *gin.Context) {
	var req ItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Status == "" {
		req.Status = "active"
	}

	database.DB.Mutex.Lock()
	defer database.DB.Mutex.Unlock()

	itemID := database.DB.GetNextID()
	item := &models.Item{
		ID:        itemID,
		Name:      req.Name,
		Status:    req.Status,
		CreatedAt: time.Now(),
	}

	database.DB.Items[itemID] = item
	c.JSON(http.StatusCreated, *item)
}

type AddToCartRequest struct {
	ItemID uint `json:"item_id" binding:"required"`
}

func AddToCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Mutex.Lock()
	defer database.DB.Mutex.Unlock()

	// Get user's cart
	var cart *models.Cart
	for _, c := range database.DB.Carts {
		if c.UserID == userID.(uint) && c.Status == "active" {
			cart = c
			break
		}
	}

	if cart == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	// Check if item exists
	item, exists := database.DB.Items[req.ItemID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	// Check if item already in cart
	cartItemKey := fmt.Sprintf("%d-%d", cart.ID, req.ItemID)
	if _, exists := database.DB.CartItems[cartItemKey]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Item already in cart"})
		return
	}

	// Add item to cart
	cartItem := &models.CartItem{
		CartID: cart.ID,
		ItemID: req.ItemID,
		Cart:   *cart,
		Item:   *item,
	}

	database.DB.CartItems[cartItemKey] = cartItem

	c.JSON(http.StatusCreated, gin.H{"message": "Item added to cart successfully"})
}

func GetCarts(c *gin.Context) {
	database.DB.Mutex.RLock()
	defer database.DB.Mutex.RUnlock()

	var carts []models.Cart
	for _, cart := range database.DB.Carts {
		cartWithItems := *cart
		cartWithItems.CartItems = []models.CartItem{}
		
		// Get cart items
		for _, cartItem := range database.DB.CartItems {
			if cartItem.CartID == cart.ID {
				cartWithItems.CartItems = append(cartWithItems.CartItems, *cartItem)
			}
		}
		
		carts = append(carts, cartWithItems)
	}

	c.JSON(http.StatusOK, carts)
}

func GetUserCart(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	database.DB.Mutex.RLock()
	defer database.DB.Mutex.RUnlock()

	// Get user's cart
	var cart *models.Cart
	for _, c := range database.DB.Carts {
		if c.UserID == userID.(uint) && c.Status == "active" {
			cart = c
			break
		}
	}

	if cart == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	cartWithItems := *cart
	cartWithItems.CartItems = []models.CartItem{}
	
	// Get cart items
	for _, cartItem := range database.DB.CartItems {
		if cartItem.CartID == cart.ID {
			cartWithItems.CartItems = append(cartWithItems.CartItems, *cartItem)
		}
	}

	c.JSON(http.StatusOK, cartWithItems)
}

func GetCartByID(c *gin.Context) {
	cartIDStr := c.Param("id")
	var cartID uint
	if _, err := fmt.Sscanf(cartIDStr, "%d", &cartID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart ID"})
		return
	}

	database.DB.Mutex.RLock()
	defer database.DB.Mutex.RUnlock()

	cart, exists := database.DB.Carts[cartID]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cart not found"})
		return
	}

	cartWithItems := *cart
	cartWithItems.CartItems = []models.CartItem{}
	
	// Get cart items
	for _, cartItem := range database.DB.CartItems {
		if cartItem.CartID == cartID {
			cartWithItems.CartItems = append(cartWithItems.CartItems, *cartItem)
		}
	}

	c.JSON(http.StatusOK, cartWithItems)
}

func CreateOrder(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	database.DB.Mutex.Lock()
	defer database.DB.Mutex.Unlock()

	// Get user's active cart
	var cart *models.Cart
	for _, c := range database.DB.Carts {
		if c.UserID == userID.(uint) && c.Status == "active" {
			cart = c
			break
		}
	}

	if cart == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active cart found"})
		return
	}

	// Check if cart has items
	hasItems := false
	for _, cartItem := range database.DB.CartItems {
		if cartItem.CartID == cart.ID {
			hasItems = true
			break
		}
	}

	if !hasItems {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Create order
	orderID := database.DB.GetNextID()
	order := &models.Order{
		ID:        orderID,
		CartID:    cart.ID,
		UserID:    cart.UserID,
		CreatedAt: time.Now(),
		Cart:      *cart,
	}

	database.DB.Orders[orderID] = order

	// Mark cart as ordered and create a new cart for the user
	cart.Status = "ordered"

	// Create new cart for user
	newCartID := database.DB.GetNextID()
	newCart := &models.Cart{
		ID:        newCartID,
		UserID:    cart.UserID,
		Name:      "Default Cart",
		Status:    "active",
		CreatedAt: time.Now(),
		CartItems: []models.CartItem{},
	}
	database.DB.Carts[newCartID] = newCart

	// Update user's cart ID
	if user, exists := database.DB.Users[cart.UserID]; exists {
		user.CartID = newCartID
	}

	c.JSON(http.StatusCreated, *order)
}

func GetOrders(c *gin.Context) {
	database.DB.Mutex.RLock()
	defer database.DB.Mutex.RUnlock()

	var orders []models.Order
	for _, order := range database.DB.Orders {
		orderWithCart := *order
		
		// Get cart with items
		if cart, exists := database.DB.Carts[order.CartID]; exists {
			cartWithItems := *cart
			cartWithItems.CartItems = []models.CartItem{}
			
			for _, cartItem := range database.DB.CartItems {
				if cartItem.CartID == cart.ID {
					cartWithItems.CartItems = append(cartWithItems.CartItems, *cartItem)
				}
			}
			
			orderWithCart.Cart = cartWithItems
		}
		
		orders = append(orders, orderWithCart)
	}

	c.JSON(http.StatusOK, orders)
}

func GetUserOrders(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	database.DB.Mutex.RLock()
	defer database.DB.Mutex.RUnlock()

	var orders []models.Order
	for _, order := range database.DB.Orders {
		if order.UserID == userID.(uint) {
			orderWithCart := *order
			
			// Get cart with items
			if cart, exists := database.DB.Carts[order.CartID]; exists {
				cartWithItems := *cart
				cartWithItems.CartItems = []models.CartItem{}
				
				for _, cartItem := range database.DB.CartItems {
					if cartItem.CartID == cart.ID {
						cartWithItems.CartItems = append(cartWithItems.CartItems, *cartItem)
					}
				}
				
				orderWithCart.Cart = cartWithItems
			}
			
			orders = append(orders, orderWithCart)
		}
	}

	c.JSON(http.StatusOK, orders)
}
