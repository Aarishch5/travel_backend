package services

import (
	"TravelBackend/database/repository"
	"TravelBackend/models"

	"github.com/jmoiron/sqlx"
)

func UpdateDriverLocation(db *sqlx.DB, driverID string, lat, lng float64) error {
	return repository.UpdateDriverLocation(db, driverID, lat, lng)
}

func GetDriverLocation(db *sqlx.DB, driverID string) (*models.DriverLocation, error) {
	return repository.DriverCurrentLocation(db, driverID)
}
