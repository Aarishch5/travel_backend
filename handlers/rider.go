package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"

	"TravelBackend/database/dbHelper"
	"TravelBackend/models"
	"TravelBackend/utils"
)

func CreateRider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only POST is allowed")
		return
	}

	var req models.CreateRiderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	// validating the email
	if msg := utils.ValidateEmail(req.Email); msg != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid email address")
	}

	// validating the phone number
	if msg := utils.ValidatePhoneNumber(req.Phone); msg != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid phone number")
	}

	exists, err := dbHelper.GetRiderByEmailOrPhone(db, req.Email, req.Phone)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusConflict, "rider with this email or phone already exists")
		return
	}

	id, err := dbHelper.CreateRider(db, req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "could not create rider")
		return
	}

	rider, err := dbHelper.GetRiderByID(db, id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "rider created but could not be fetched")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "rider created successfully",
		"rider":   rider,
	})
}

func DeleteRider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodDelete {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only DELETE is allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "id query parameter is required")
		return
	}

	err := dbHelper.DeleteRider(db, id)
	if err == models.ErrRiderNotFound {
		utils.RespondError(w, http.StatusNotFound, "rider not found")
		return
	}
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "could not delete rider")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "rider deleted successfully",
	})
}
