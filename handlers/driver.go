package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"

	"TravelBackend/database/dbHelper"
	"TravelBackend/models"
	"TravelBackend/utils"
)

func CreateDriver(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only POST is allowed")
		return
	}

	var req models.CreateDriverRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// validating the name
	if strings.TrimSpace(req.Name) == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	// Validating the email
	if msg := utils.ValidateEmail(req.Email); msg != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid email address")
	}

	// validating the phone number
	if msg := utils.ValidatePhoneNumber(req.Phone); msg != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid phone number")
	}

	// validate the license
	if strings.TrimSpace(req.LicenseNumber) == "" {
		http.Error(w, "license_number is required", http.StatusBadRequest)
		return
	}

	exists, err := dbHelper.GetDriverByEmailOrPhone(db, req.Email, req.Phone)
	if err != nil {
		log.Println("GetDriverByEmailOrPhone error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	if exists {
		utils.RespondError(w, http.StatusConflict, "driver with this email or phone already exists")
		return
	}

	id, err := dbHelper.CreateDriver(db, req)
	if err != nil {
		log.Println("CreateDriver error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not create driver")
		return
	}

	driver, err := dbHelper.GetDriverByID(db, id)
	if err != nil {
		log.Println("GetDriverByID error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "driver created but could not be fetched")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, map[string]interface{}{
		"message": "driver created successfully",
		"driver":  driver,
	})
}

func DeleteDriver(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodDelete {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only DELETE is allowed")
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		utils.RespondError(w, http.StatusBadRequest, "id query parameter is required")
		return
	}

	err := dbHelper.DeleteDriver(db, id)
	if err == models.ErrDriverNotFound {
		utils.RespondError(w, http.StatusNotFound, "driver not found")
		return
	}
	if err != nil {
		log.Println("DeleteDriver error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not delete driver")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "driver deleted successfully",
	})
}
