package repository

import "github.com/jmoiron/sqlx"

func CreateRideOffers(db *sqlx.DB, rideID string, driverIDs []string) error {
	if len(driverIDs) == 0 {
		return nil
	}

	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, driverID := range driverIDs {
		if _, err := tx.Exec(`
			INSERT INTO ride_offers (ride_id, driver_id, status)
			VALUES ($1, $2, 'PENDING')`, rideID, driverID); err != nil {
			return err
		}
	}

	return tx.Commit()
}
