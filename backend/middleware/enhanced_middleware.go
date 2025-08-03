package middleware

import (
	"ecommerce-backend/database"
	"ecommerce-backend/models"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Rate limiter using in-memory store
type RateLimiter struct {
	requests map[string][]time.Time
	mutex    sync.RWMutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	now := time.Now()
	
	// Clean old requests
	if requests, exists := rl.requests[ip]; exists {
		var validRequests []time.Time
		for _, requestTime := range requests {
			if now.Sub(requestTime) < rl.window {
				validRequests = append(validRequests, requestTime)
			}
		}
		rl.requests[ip] = validRequests
	}

	// Check if limit exceeded
	if len(rl.requests[ip]) >= rl.limit {
		return false
	}

	// Add current request
	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

var globalRateLimiter *RateLimiter

// RateLimit middleware
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	if globalRateLimiter == nil {
		globalRateLimiter = NewRateLimiter(limit, window)
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		
		if !globalRateLimiter.Allow(clientIP) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "Rate limit exceeded. Please try again later.",
				"meta": gin.H{
					"timestamp": time.Now().Format(time.RFC3339),
					"version":   "v1.0",
					"retry_after": window.Seconds(),
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// Security headers middleware
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Security headers
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:")
		
		// HSTS (only enable in production with HTTPS)
		// c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		
		c.Next()
	}
}

// Request ID middleware
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.New().String()
		c.Header("X-Request-ID", requestID)
		c.Set("RequestID", requestID)
		c.Next()
	}
}

// Get user profile endpoint
func GetUserProfile() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "User not authenticated",
				"meta": gin.H{
					"timestamp": time.Now().Format(time.RFC3339),
					"version":   "v1.0",
				},
			})
			return
		}

		var user models.User
		if err := database.DB.Select("id, username, created_at, updated_at").First(&user, userID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "User not found",
				"meta": gin.H{
					"timestamp": time.Now().Format(time.RFC3339),
					"version":   "v1.0",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "User profile retrieved successfully",
			"data": gin.H{
				"id":         user.ID,
				"username":   user.Username,
				"role":       "admin",
				"created_at": user.CreatedAt,
				"updated_at": user.UpdatedAt,
			},
			"meta": gin.H{
				"timestamp": time.Now().Format(time.RFC3339),
				"version":   "v1.0",
			},
		})
	}
}

// Clear cart endpoint
func ClearCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"error":   "User not authenticated",
				"meta": gin.H{
					"timestamp": time.Now().Format(time.RFC3339),
					"version":   "v1.0",
				},
			})
			return
		}

		// Find user's cart
		var cart models.Cart
		if err := database.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"error":   "Cart not found",
				"meta": gin.H{
					"timestamp": time.Now().Format(time.RFC3339),
					"version":   "v1.0",
				},
			})
			return
		}

		// Clear all cart items
		if err := database.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to clear cart",
				"meta": gin.H{
					"timestamp": time.Now().Format(time.RFC3339),
					"version":   "v1.0",
				},
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Cart cleared successfully",
			"data": gin.H{
				"cart_id": cart.ID,
				"items_removed": "all",
			},
			"meta": gin.H{
				"timestamp": time.Now().Format(time.RFC3339),
				"version":   "v1.0",
			},
		})
	}
}
