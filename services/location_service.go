package services

import (
	"TravelBackend/database/dbHelper"

	"github.com/jmoiron/sqlx"
)

func UpdateDriverLocation(db *sqlx.DB, driverID string, lat, lng float64) error {
	return dbHelper.UpsertDriverLocation(db, driverID, lat, lng)
}
