package handlers

import (
	"encoding/json"
	"log"
	"net/http"

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

	var req models.UpdateLocationRequest
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
		"message": "location updated successfully via PostGIS",
	})
}
