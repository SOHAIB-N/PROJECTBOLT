package main

import (
	"log"
	"net/http"
	"os"
	"restaurant-system/handlers"
	"restaurant-system/middleware"
	"restaurant-system/models"
	"restaurant-system/utils"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&models.User{}, &models.Payment{}, &models.MenuItem{}, &models.Order{}, &models.OrderItem{}, &models.Admin{})

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	menuHandler := handlers.NewMenuHandler(db)
	orderHandler := handlers.NewOrderHandler(db)
	adminHandler := handlers.NewAdminHandler(db)

	// Initialize router
	r := mux.NewRouter()

	// CORS middleware
	r.Use(middleware.CORSMiddleware)

	// Public routes
	r.HandleFunc("/api/auth/register", authHandler.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/auth/login", authHandler.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/menu", menuHandler.GetMenu).Methods("GET", "OPTIONS")

	// Protected routes
	api := r.PathPrefix("/api").Subrouter()
	api.Use(middleware.AuthMiddleware)

	// User routes
	api.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST", "OPTIONS")
	api.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET", "OPTIONS")
	api.HandleFunc("/orders/user", orderHandler.GetUserOrders).Methods("GET", "OPTIONS")

	// Admin routes
	admin := r.PathPrefix("/api/admin").Subrouter()
	admin.Use(middleware.AdminAuthMiddleware)

	admin.HandleFunc("/menu", menuHandler.AddMenuItem).Methods("POST", "OPTIONS")
	admin.HandleFunc("/menu/{id}", menuHandler.UpdateMenuItem).Methods("PUT", "OPTIONS")
	admin.HandleFunc("/menu/{id}", menuHandler.DeleteMenuItem).Methods("DELETE", "OPTIONS")
	admin.HandleFunc("/orders", adminHandler.GetAllOrders).Methods("GET", "OPTIONS")
	admin.HandleFunc("/orders/{id}/status", orderHandler.UpdateOrderStatus).Methods("PUT", "OPTIONS")

	// WebSocket endpoint for real-time updates
	r.HandleFunc("/ws", utils.HandleWebSocket)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}