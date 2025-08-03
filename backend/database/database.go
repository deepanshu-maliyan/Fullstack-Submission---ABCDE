package database

import (
	"ecommerce-backend/models"
	"log"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// In-memory database using maps
type InMemoryDB struct {
	Users     map[uint]*models.User
	Items     map[uint]*models.Item
	Carts     map[uint]*models.Cart
	CartItems map[string]*models.CartItem // key: "cartID-itemID"
	Orders    map[uint]*models.Order
	Mutex     sync.RWMutex
	nextID    uint
}

var DB *InMemoryDB

func Connect() {
	DB = &InMemoryDB{
		Users:     make(map[uint]*models.User),
		Items:     make(map[uint]*models.Item),
		Carts:     make(map[uint]*models.Cart),
		CartItems: make(map[string]*models.CartItem),
		Orders:    make(map[uint]*models.Order),
		nextID:    1,
	}

	// Seed some initial items
	seedItems()
	// Create admin user
	seedAdminUser()
	log.Println("In-memory database initialized successfully")
}

func (db *InMemoryDB) GetNextID() uint {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()
	id := db.nextID
	db.nextID++
	return id
}

func seedItems() {
	if len(DB.Items) == 0 {
		items := []models.Item{
			{ID: DB.GetNextID(), Name: "Laptop", Status: "active", Image: "/assets/products/laptop.jpg", CreatedAt: time.Now()},
			{ID: DB.GetNextID(), Name: "Smartphone", Status: "active", Image: "/assets/products/smartphone.jpg", CreatedAt: time.Now()},
			{ID: DB.GetNextID(), Name: "Headphones", Status: "active", Image: "/assets/products/headphones.jpg", CreatedAt: time.Now()},
			{ID: DB.GetNextID(), Name: "Keyboard", Status: "active", Image: "/assets/products/keyboard.jpg", CreatedAt: time.Now()},
			{ID: DB.GetNextID(), Name: "Mouse", Status: "active", Image: "/assets/products/mouse.jpg", CreatedAt: time.Now()},
			{ID: DB.GetNextID(), Name: "Monitor", Status: "active", Image: "/assets/products/monitor.jpg", CreatedAt: time.Now()},
			{ID: DB.GetNextID(), Name: "Tablet", Status: "active", Image: "/assets/products/tablet.jpg", CreatedAt: time.Now()},
			{ID: DB.GetNextID(), Name: "Webcam", Status: "active", Image: "/assets/products/webcam.jpg", CreatedAt: time.Now()},
		}

		for _, item := range items {
			DB.Items[item.ID] = &item
		}
		log.Println("Seeded initial items with image URLs")
	}
}

func seedAdminUser() {
	// Check if admin user already exists
	for _, user := range DB.Users {
		if user.Username == "admin" {
			log.Println("Admin user already exists")
			return
		}
	}

	// Hash the admin password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Admin@123"), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error creating admin user: %v", err)
		return
	}

	// Create admin user
	adminID := DB.GetNextID()
	adminUser := &models.User{
		ID:        adminID,
		Username:  "admin",
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
	}

	// Create a cart for the admin user
	cartID := DB.GetNextID()
	adminCart := &models.Cart{
		ID:        cartID,
		UserID:    adminID,
		Name:      "Admin Cart",
		Status:    "active",
		CreatedAt: time.Now(),
		CartItems: []models.CartItem{},
	}

	// Store in database
	DB.Carts[cartID] = adminCart
	adminUser.CartID = cartID
	DB.Users[adminID] = adminUser

	log.Println("Created admin user (username: admin, password: Admin@123)")
}
