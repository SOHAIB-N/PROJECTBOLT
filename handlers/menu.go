package handlers

import (
	"encoding/json"
	"net/http"
	"restaurant-system/models"
	"restaurant-system/utils"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type MenuHandler struct {
	DB *gorm.DB
}

func NewMenuHandler(db *gorm.DB) *MenuHandler {
	return &MenuHandler{DB: db}
}

func (h *MenuHandler) GetMenu(w http.ResponseWriter, r *http.Request) {
	var menuItems []models.MenuItem
	if err := h.DB.Where("available = ?", true).Find(&menuItems).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error fetching menu")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, menuItems)
}

func (h *MenuHandler) AddMenuItem(w http.ResponseWriter, r *http.Request) {
	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if item.Name == "" || item.Price <= 0 {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid menu item details")
		return
	}

	if err := h.DB.Create(&item).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error creating menu item")
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, item)
}

func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var item models.MenuItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.DB.Model(&models.MenuItem{}).Where("id = ?", id).Updates(item).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error updating menu item")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, item)
}

func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.DB.Delete(&models.MenuItem{}, id).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Error deleting menu item")
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "Menu item deleted successfully"})
}