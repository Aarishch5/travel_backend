package dbHelper

import (
	"github.com/jmoiron/sqlx"
)

func UpsertDriverLocation(db *sqlx.DB, driverID string, lat, lng float64) error {
	query := `
		INSERT INTO driver_locations (driver_id, location, updated_at)
		VALUES ($1, ST_MakePoint($2, $3)::geography, now())
		ON CONFLICT (driver_id)
		DO UPDATE SET location = ST_MakePoint($4, $5)::geography, updated_at = now()`

	_, err := db.Exec(query, driverID, lng, lat, lng, lat)
	return err
}
