package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"TravelBackend/database/dbHelper"
	"TravelBackend/middleware"
	"TravelBackend/models"
	"TravelBackend/services"
	"TravelBackend/utils"
)

func RequestRide(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only POST is allowed")
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.RequestRideRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PickupLat == 0 && req.PickupLng == 0 {
		utils.RespondError(w, http.StatusBadRequest, "pickup_lat and pickup_lng are required")
		return
	}
	if req.DropLat == 0 && req.DropLng == 0 {
		utils.RespondError(w, http.StatusBadRequest, "dropping latitude and longitude are required")
		return
	}

	ride, err := services.RequestRide(db, claims.UserID, req)
	if err != nil {
		log.Println("RequestRide error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not request ride")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, ride)
}

func AcceptRide(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPatch {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only PATCH is allowed")
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	rideID := r.PathValue("id")
	if rideID == "" {
		utils.RespondError(w, http.StatusBadRequest, "ride id is required")
		return
	}

	ride, err := services.AcceptRide(db, rideID, claims.UserID)
	switch err {
	case nil:
		utils.RespondJSON(w, http.StatusOK, ride)
	case models.ErrRideAlreadyTaken:
		utils.RespondError(w, http.StatusConflict, "another driver already accepted this ride")
	case models.ErrOfferNotFound:
		utils.RespondError(w, http.StatusConflict, "you were not offered this ride or already responded")
	case models.ErrRideNotFound:
		utils.RespondError(w, http.StatusNotFound, "ride not found")
	default:
		log.Println("AcceptRide error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not accept ride")
	}
}

func RejectRide(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPatch {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only PATCH is allowed")
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	rideID := r.PathValue("id")
	if rideID == "" {
		utils.RespondError(w, http.StatusBadRequest, "ride id is required")
		return
	}

	ride, err := services.RejectRide(db, rideID, claims.UserID)
	switch err {
	case nil:
		utils.RespondJSON(w, http.StatusOK, ride)
	case models.ErrOfferNotFound:
		utils.RespondError(w, http.StatusConflict, "you were not offered this ride or already responded")
	case models.ErrRideNotFound:
		utils.RespondError(w, http.StatusNotFound, "ride not found")
	default:
		log.Println("RejectRide error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not reject ride")
	}
}

func GetPendingRides(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only GET is allowed")
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	rides, err := dbHelper.GetPendingRidesForDriver(db, claims.UserID)
	if err != nil {
		log.Println("GetPendingRides error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not fetch pending rides")
		return
	}

	utils.RespondJSON(w, http.StatusOK, rides)
}

func RideCompleted(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodPatch {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only PATCH is allowed")
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	rideID := r.PathValue("id")
	if rideID == "" {
		utils.RespondError(w, http.StatusBadRequest, "ride id is required")
		return
	}

	ride, err := services.CompleteRide(db, rideID, claims.UserID)
	switch err {
	case nil:
		utils.RespondJSON(w, http.StatusOK, ride)
	case models.ErrRideNotActive:
		utils.RespondError(w, http.StatusConflict, "ride is not assigned or  not accepted")
	case models.ErrNotAtDroppingLocation:
		utils.RespondError(w, http.StatusConflict, "you are too far from the drop location points")
	case models.ErrRideNotFound:
		utils.RespondError(w, http.StatusNotFound, "ride not found")
	default:
		log.Println("RideCompleted error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not complete ride")
	}

}

//func CalculateFair(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
//
//}
