package repository

import (
	"TravelBackend/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func UpdateDriverLocation(db *sqlx.DB, driverID string, lat, lng float64) error {
	//query1 := `SELECT COUNT(*) FROM driver_locations WHERE driver_id = $1;`

	query := `
		INSERT INTO driver_locations (driver_id, latitude, longitude, updated_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (driver_id)
		DO UPDATE SET latitude = $2, longitude = $3, updated_at = now()`

	_, err := db.Exec(query, driverID, lat, lng)

	return err
}

func DriverCurrentLocation(db *sqlx.DB, driver_id string) (*models.DriverLocation, error) {

	query := `SELECT driver_id, latitude, longitude, updated_at FROM driver_locations WHERE driver_id = $1`

	var driverLocationInfo models.DriverLocation

	err := db.Get(&driverLocationInfo, query, driver_id)

	if err == sql.ErrNoRows {
		return nil, models.ErrDriverNotFound
	}
	if err != nil {
		return nil, err
	}

	return &driverLocationInfo, nil
}
