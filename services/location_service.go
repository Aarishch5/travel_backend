package services

import (
	"TravelBackend/database/dbHelper"
	"TravelBackend/models"

	"github.com/jmoiron/sqlx"
)

func UpdateDriverLocation(db *sqlx.DB, driverID string, lat, lng float64) error {
	return dbHelper.UpsertDriverLocation(db, driverID, lat, lng)
}

func DriverLocation(db *sqlx.DB, driverID string) (*models.DriverLocation, error) {
	return dbHelper.DriverCurrentLocation(db, driverID)
}
