package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"

	"TravelBackend/middleware"
	"TravelBackend/models"
	"TravelBackend/services"
	"TravelBackend/utils"
)

func UpdateDriverLocation(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req models.DriverLocation
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		utils.RespondError(w, http.StatusBadRequest, "latitude must be between -90 and 90")
		return
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		utils.RespondError(w, http.StatusBadRequest, "longitude must be between -180 and 180")
		return
	}

	if err := services.UpdateDriverLocation(db, claims.UserID, req.Latitude, req.Longitude); err != nil {
		log.Println("UpdateDriverLocation error:", err)
		utils.RespondError(w, http.StatusInternalServerError, "could not update location")
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{
		"message": "location updated successfully",
	})
}

func GetDriverLocation(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	if r.Method != http.MethodGet {
		utils.RespondError(w, http.StatusMethodNotAllowed, "only GET is allowed")
		return
	}

	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	driverID := claims.UserID
	if driverID == "" || strings.TrimSpace(driverID) == "" {
		utils.RespondError(w, http.StatusNotFound, "driver id or user is required")
		return
	}

	driverLocation, err := services.DriverLocation(db, driverID)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "driver location not found")
		return
	}
	utils.RespondJSON(w, http.StatusOK, driverLocation)
}

func DriverLocationHandler(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {
	switch r.Method {
	case http.MethodGet:
		GetDriverLocation(w, r, db)
	case http.MethodPatch:
		UpdateDriverLocation(w, r, db)
	default:
		utils.RespondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
