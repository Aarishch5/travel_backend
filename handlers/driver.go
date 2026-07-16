package handlers

import (
	"TravelBackend/middleware"
	"TravelBackend/services"
	"log"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"

	"TravelBackend/database/repository"
	"TravelBackend/models"
	"TravelBackend/utils"
)

func RegisterDriver(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only POST is allowed")
		return
	}

	var req models.CreateDriverRequest
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	utils.RespondError(w, http.StatusBadRequest, "invalid request body")
	//	return
	//}

	if err := utils.ParseBody(r.Body, &req); err != nil {
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
		return
	}

	// validating the phone number
	if msg := utils.ValidatePhoneNumber(req.Phone); msg != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid phone number")
		return
	}

	// validate the license
	if strings.TrimSpace(req.LicenseNumber) == "" {
		http.Error(w, "license_number is required", http.StatusBadRequest)
		return
	}

	// Validate the password
	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	driver, err := services.RegisterDriver(db, req)
	if err == models.ErrEmailOrPhoneExists {
		utils.RespondError(w, http.StatusConflict, "driver with this email or phone already exists")
		return
	}

	if err != nil {
		log.Println("RegisterDriver error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not register driver")
		return
	}

	token, err := utils.GenerateToken(driver.ID, "driver")
	if err != nil {
		log.Println("GenerateToken error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "driver registered but token generation failed")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  driver,
	})
}

func UpdateDriverStatus(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPatch {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only PATCH is allowed")
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.UpdateDriverStatusRequest
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	utils.RespondError(w, http.StatusBadRequest, "invalid request body")
	//	return
	//}
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if !models.IsValidDriverStatus(req.Status) {
		utils.RespondError(w, http.StatusBadRequest, "status must be one of ONLINE, OFFLINE, ON_TRIP")
		return
	}

	err := services.UpdateDriverStatus(db, claims.UserID, req.Status)
	if err == models.ErrDriverNotFound {
		utils.RespondError(w, http.StatusNotFound, "driver not found")
		return
	}
	if err != nil {
		log.Println("UpdateDriverStatus error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not update driver status")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "driver status updated successfully",
		"status":  req.Status,
	})
}

func LoginDriver(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only POST is allowed")
		return
	}

	var req models.LoginDriverRequest
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	utils.RespondError(w, http.StatusBadRequest, "invalid request body")
	//	return
	//}
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, driver, err := services.LoginDriver(db, req)
	if err == models.ErrInvalidCredentials {
		utils.RespondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		log.Println("LoginDriver error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	utils.RespondJSON(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  driver,
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

	err := repository.DeleteDriver(db, id)
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
