package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restaurant-system/models"
	"restaurant-system/utils"
	"time"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type OrderHandler struct {
	DB *gorm.DB
}

func NewOrderHandler(db *gorm.DB) *OrderHandler {
	return &OrderHandler{DB: db}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userID := utils.GetUserIDFromContext(r.Context())
	order.UserID = userID
	order.Status = "pending"
	order.OrderNumber = utils.GenerateOrderNumber()

	// Calculate total amount and validate items
	var totalAmount float64
	for i, item := range order.Items {
		var menuItem models.MenuItem
		if err := h.DB.First(&menuItem, item.MenuItemID).Error; err != nil {
			utils.RespondWithError(w, http.StatusBadRequest, "Invalid menu item")
			return
		}
		order.Items[i].Price = menuItem.Price
		order.Items[i].TotalPrice = menuItem.Price * float64(item.Quantity)
		totalAmount += order.Items[i].TotalPrice
	}
	order.TotalAmount = totalAmount

	// Validate delivery time if it's a hostel delivery
	if order.IsDelivery {
		now := time.Now()
		hour := now.Hour()
		if hour < 19 || hour >= 22 {
			utils.RespondWithError(w, http.StatusBadRequest, "Hostel delivery is only available between 7 PM and 10 PM")
			return
		}
		if order.RoomNumber == "" {
			utils.RespondWithError(w, http.StatusBadRequest, "Room number is required for hostel delivery")
			return
		}
	}

	if err := h.DB.Create(&order).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating order")
		return
	}

	// Notify kitchen staff via WebSocket
	utils.BroadcastNewOrder(order)

	utils.RespondWithJSON(w, http.StatusCreated, order)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var order models.Order
	if err := h.DB.Preload("Items.MenuItem").First(&order, orderID).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	userID := utils.GetUserIDFromContext(r.Context())
	if order.UserID != userID {
		utils.RespondWithError(w, http.StatusForbidden, "Unauthorized access to order")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, order)
}

func (h *OrderHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	userID := utils.GetUserIDFromContext(r.Context())

	var orders []models.Order
	if err := h.DB.Where("user_id = ?", userID).Preload("Items.MenuItem").Find(&orders).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching orders")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	var statusUpdate struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var order models.Order
	if err := h.DB.First(&order, orderID).Error; err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	if err := h.DB.Model(&order).Update("status", statusUpdate.Status).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating order status")
		return
	}

	// Print receipt if order is confirmed
	if statusUpdate.Status == "confirmed" {
		if err := utils.PrintReceipt(order); err != nil {
			fmt.Printf("Error printing receipt: %v\n", err)
		}
	}

	utils.RespondWithJSON(w, http.StatusOK, order)
}