package handlers

import (
	"net/http"
	"restaurant-system/models"
	"restaurant-system/utils"

	"gorm.io/gorm"
)

type AdminHandler struct {
	DB *gorm.DB
}

func NewAdminHandler(db *gorm.DB) *AdminHandler {
	return &AdminHandler{DB: db}
}

func (h *AdminHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	var orders []models.Order
	if err := h.DB.Preload("Items.MenuItem").Preload("User").Find(&orders).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching orders")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, orders)
}

func (h *AdminHandler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	var totalOrders int64
	var totalRevenue float64
	var pendingOrders int64

	h.DB.Model(&models.Order{}).Count(&totalOrders)
	h.DB.Model(&models.Order{}).Where("status = ?", "pending").Count(&pendingOrders)
	h.DB.Model(&models.Order{}).Select("COALESCE(SUM(total_amount), 0)").Scan(&totalRevenue)

	stats := map[string]interface{}{
		"totalOrders":    totalOrders,
		"totalRevenue":   totalRevenue,
		"pendingOrders":  pendingOrders,
	}

	utils.RespondWithJSON(w, http.StatusOK, stats)
}

func (h *AdminHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	var menuItem models.MenuItem
	if err := utils.ParseBody(r, &menuItem); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.DB.Create(&menuItem).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating menu item")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, menuItem)
}