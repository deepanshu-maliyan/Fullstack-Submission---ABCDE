
# E-commerce Full Stack Application

A complete e-commerce web application with Go backend and React frontend.

## Project Structure

```
Assessment - FullStack(ABCDE)/
├── backend/                 # Go backend with Gin and GORM
│   ├── handlers/           # API handlers
│   ├── models/            # Database models
│   ├── middleware/        # Auth middleware
│   ├── database/          # Database connection and setup
│   ├── utils/             # JWT utilities
│   ├── main.go           # Main application file
│   ├── go.mod            # Go dependencies
│   └── README.md         # Backend documentation
├── frontend/              # React frontend
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── AuthContext.js # Authentication context
│   │   ├── api.js        # API client
│   │   └── App.js        # Main app component
│   ├── package.json      # Node dependencies
│   └── public/           # Static files
└── README.md             # This file
```

## Features

### Backend (Go + Gin + In-Memory DB)
- Admin authentication with pre-configured credentials
- Shopping cart management
- Order processing
- RESTful API design
- In-memory database with sample data
- Comprehensive test suite with Ginkgo
- CORS enabled for frontend integration

### Frontend (React)
- Admin login interface
- Items listing and cart management
- Real-time notifications with toast messages
- Responsive design
- Protected routes with authentication
- Order history and cart viewing

## Quick Start

### Prerequisites
- Go 1.21 or higher
- Node.js 16 or higher
- npm or yarn

### Backend Setup

1. Navigate to the backend directory:
   ```cmd
   cd backend
   ```

2. Install Go dependencies:
   ```cmd
   go mod tidy
   ```

3. Run the backend server:
   ```cmd
   go run main.go
   ```

The backend will start on `http://localhost:8080`

### Frontend Setup

1. Open a new terminal and navigate to the frontend directory:
   ```cmd
   cd frontend
   ```

2. Install Node dependencies:
   ```cmd
   npm install
   ```

3. Start the React development server:
   ```cmd
   npm start
   ```

The frontend will start on `http://localhost:3000`

## API Endpoints

### Public Endpoints
- `POST /users` - Create a new user
- `POST /users/login` - User login
- `GET /items` - List all items

### Protected Endpoints (require Authorization header)
- `GET /users` - List all users
- `POST /items` - Create a new item
- `POST /carts` - Add item to cart
- `GET /carts/user` - Get current user's cart
- `GET /carts` - List all carts
- `POST /orders` - Create order from cart
- `GET /orders/user` - Get current user's orders
- `GET /orders` - List all orders

## Usage Flow

1. **Registration/Login**: 
   - New users can register with username/password
   - Existing users can login to get access token

2. **Shopping**:
   - Browse available items on the main page
   - Click items to add them to cart
   - View cart contents with the "Cart" button

3. **Checkout**:
   - Click "Checkout" to convert cart to order
   - View order history with "Order History" button

## Testing

### Backend Tests
```cmd
cd backend
go test ./...
```

### Frontend Tests
```cmd
cd frontend
npm test
```

## Database Schema

The application uses these entities:
- **Users**: Authentication and user management
- **Items**: Product catalog
- **Carts**: User shopping carts
- **CartItems**: Items in carts (many-to-many)
- **Orders**: Completed purchases

## Technologies Used

### Backend
- Go 1.21
- Gin (Web framework)
- GORM (ORM)
- JWT for authentication
- SQLite (Database)
- Ginkgo (Testing)
- bcrypt (Password hashing)

### Frontend
- React 18
- React Router (Navigation)
- Axios (HTTP client)
- React Toastify (Notifications)
- CSS3 (Styling)

## Development Notes

- The backend automatically seeds sample items on first run
- JWT tokens are valid for 24 hours
- Users can only have one active cart at a time
- Orders are created from cart contents and create a new empty cart
- The frontend has proper error handling and user feedback
- All routes are protected appropriately with authentication

## Production Considerations

For production deployment, consider:
- Using a production database (PostgreSQL, MySQL)
- Environment-based configuration
- HTTPS/TLS encryption
- Rate limiting and security headers
- Logging and monitoring
- Docker containerization
- CI/CD pipeline setup
