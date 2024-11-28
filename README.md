# Food Court Backend System

A robust backend system for a food court management system built with Go.

## Features

- User Authentication (JWT)
- Menu Management
- Order Processing
- Real-time Updates (WebSocket)
- Payment Integration
- Thermal Printer Support
- Admin Dashboard

## Setup

1. Install Go 1.21 or later
2. Clone the repository
3. Copy `.env.example` to `.env` and configure your environment variables
4. Install dependencies:
   ```bash
   go mod download
   ```
5. Start the server:
   ```bash
   go run main.go
   ```

## API Endpoints

### Authentication
- POST /api/auth/register - Register new user
- POST /api/auth/login - User login

### Menu
- GET /api/menu - Get menu items
- POST /api/menu - Add menu item (Admin)
- PUT /api/menu/{id} - Update menu item (Admin)
- DELETE /api/menu/{id} - Delete menu item (Admin)

### Orders
- POST /api/orders - Create new order
- GET /api/orders/{id} - Get order details
- PUT /api/orders/{id}/status - Update order status

### WebSocket
- WS /ws - Real-time order updates

## Security

- JWT-based authentication
- Password hashing with bcrypt
- Input validation and sanitization
- Rate limiting
- CORS protection

## Database Schema

The system uses PostgreSQL with the following main tables:
- users
- menu_items
- orders
- order_items
- payments
- admins

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a new Pull Request