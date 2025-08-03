# E-commerce Backend API

A simple e-commerce backend built with Go, Gin, and GORM.

## Features

- Admin user authentication (pre-configured)
- JWT-based authorization
- Shopping cart management
- Order processing
- RESTful API design
- In-memory database with sample data

## Installation

1. Make sure you have Go 1.21+ installed
2. Navigate to the backend directory
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run the application:
   ```bash
   go run main.go
   ```

The server will start on `http://localhost:8080`

## API Endpoints

### Public Endpoints

- `POST /users/login` - Admin login (username: admin, password: Admin@123)
- `GET /items` - List all items

### Protected Endpoints (require Authorization header with Bearer token)

#### Users
- `GET /users` - List all users

#### Items
- `POST /items` - Create a new item

#### Carts
- `POST /carts` - Add item to cart
- `GET /carts` - List all carts
- `GET /carts/user` - Get current user's cart
- `GET /carts/:id` - Get cart by ID

#### Orders
- `POST /orders` - Create order from cart
- `GET /orders` - List all orders
- `GET /orders/user` - Get current user's orders

## Testing

Run tests using Ginkgo:

```bash
go test ./...
```

## Database Schema

The application uses the following entities:
- Users (with authentication)
- Items (products)
- Carts (user shopping carts)
- CartItems (items in carts)
- Orders (completed purchases)

## Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```
