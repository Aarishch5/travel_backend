package dbHelper

import (
	"github.com/jmoiron/sqlx"
)

func UpsertDriverLocation(db *sqlx.DB, driverID string, lat, lng float64) error {
	query := `
		INSERT INTO driver_locations (driver_id, latitude, longitude, updated_at)
		VALUES ($1, $2, $3, now())
		ON CONFLICT (driver_id)
		DO UPDATE SET latitude = $2, longitude = $3, updated_at = now()`
	_, err := db.Exec(query, driverID, lat, lng)
	return err
}
