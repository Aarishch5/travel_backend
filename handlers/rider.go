package handlers

import (
	"TravelBackend/services"
	"log"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"

	"TravelBackend/database/dbHelper"
	"TravelBackend/models"
	"TravelBackend/utils"
)

func RegisterRider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only POST is allowed")
		return
	}

	var req models.CreateRiderRequest
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	utils.RespondError(w, http.StatusBadRequest, "invalid request body")
	//	return
	//}
	if err := utils.ParseBody(r.Body, &req); err != nil {
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
		return
	}

	// validating the phone number
	if msg := utils.ValidatePhoneNumber(req.Phone); msg != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid phone number")
		return
	}

	// Validate the password
	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	rider, err := services.RegisterRider(db, req)
	if err == models.ErrEmailOrPhoneExists {
		utils.RespondError(w, http.StatusConflict, "rider with this email or phone already exists")
		return
	}

	if err != nil {
		log.Println("RegisterRider error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not register rider")
		return
	}

	token, err := utils.GenerateToken(rider.ID, "rider")
	if err != nil {
		log.Println("GenerateToken error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "rider registered but token generation failed")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, models.AuthResponse{
		Token: token,
		User:  rider,
	})

}

func LoginRider(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only POST is allowed")
		return
	}

	var req models.LoginRiderRequest
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	utils.RespondError(w, http.StatusBadRequest, "invalid request body")
	//	return
	//}
	if err := utils.ParseBody(r.Body, &req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, rider, err := services.LoginRider(db, req)
	if err == models.ErrInvalidCredentials {
		utils.RespondError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}
	if err != nil {
		log.Println("LoginRider error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	utils.RespondJSON(w, http.StatusOK, models.AuthResponse{
		Token: token,
		User:  rider,
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
