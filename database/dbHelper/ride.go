package dbHelper

import (
	"database/sql"

	"github.com/jmoiron/sqlx"

	"TravelBackend/models"
)

func FindNearbyDrivers(db *sqlx.DB, lat, lng float64) ([]models.NearbyDriver, error) {
	query := `
		SELECT driver_id, distance_km FROM (
			SELECT dl.driver_id,
			       (6371 * acos(
			           cos(radians($1)) * cos(radians(dl.latitude)) * cos(radians(dl.longitude) - radians($2))
			           + sin(radians($1)) * sin(radians(dl.latitude))
			       )) AS distance_km
			FROM driver_locations dl
			JOIN drivers d ON d.id = dl.driver_id
			WHERE d.status = 'ONLINE'
		) sub
		WHERE distance_km < 5
		ORDER BY distance_km
		LIMIT 5`

	var drivers []models.NearbyDriver
	err := db.Select(&drivers, query, lat, lng)
	return drivers, err
}

func CreateRide(db *sqlx.DB, riderID string, req models.RequestRideRequest) (string, error) {
	var id string
	query := `
		INSERT INTO rides (rider_id, pickup_lat, pickup_lng, drop_lat, drop_lng, status)
		VALUES ($1, $2, $3, $4, $5, 'REQUESTED')
		RETURNING id`
	err := db.Get(&id, query, riderID, req.PickupLat, req.PickupLng, req.DropLat, req.DropLng)
	return id, err
}

func MarkRideNoDriversFound(db *sqlx.DB, rideID string) error {
	_, err := db.Exec(`
		UPDATE rides SET status = 'NO_DRIVERS_FOUND', updated_at = now()
		WHERE id = $1`, rideID)
	return err
}

func GetRideByID(db *sqlx.DB, rideID string) (*models.Ride, error) {
	var ride models.Ride
	query := `
		SELECT id, rider_id, driver_id, status, pickup_lat, pickup_lng, drop_lat, drop_lng,
		       requested_at, accepted_at, updated_at
		FROM rides WHERE id = $1`
	err := db.Get(&ride, query, rideID)
	if err == sql.ErrNoRows {
		return nil, models.ErrRideNotFound
	}
	if err != nil {
		return nil, err
	}
	return &ride, nil
}

func AcceptRide(db *sqlx.DB, rideID, driverID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var status string
	err = tx.Get(&status, `SELECT status FROM rides WHERE id = $1 FOR UPDATE`, rideID)
	if err == sql.ErrNoRows {
		return models.ErrRideNotFound
	}
	if err != nil {
		return err
	}
	if status != models.RideStatusRequested {
		return models.ErrRideAlreadyTaken
	}

	result, err := tx.Exec(`
		UPDATE ride_offers SET status = 'ACCEPTED', updated_at = now()
		WHERE ride_id = $1 AND driver_id = $2 AND status = 'PENDING'`,
		rideID, driverID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrOfferNotFound
	}

	if _, err := tx.Exec(`
		UPDATE rides SET driver_id = $1, status = 'ACCEPTED', accepted_at = now(), updated_at = now()
		WHERE id = $2`, driverID, rideID); err != nil {
		return err
	}

	if _, err := tx.Exec(`
		UPDATE ride_offers SET status = 'EXPIRED', updated_at = now()
		WHERE ride_id = $1 AND driver_id != $2 AND status = 'PENDING'`,
		rideID, driverID); err != nil {
		return err
	}

	return tx.Commit()
}

func RejectRide(db *sqlx.DB, rideID, driverID string) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.Exec(`
		UPDATE ride_offers SET status = 'REJECTED', updated_at = now()
		WHERE ride_id = $1 AND driver_id = $2 AND status = 'PENDING'`,
		rideID, driverID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrOfferNotFound
	}

	var remaining int
	err = tx.Get(&remaining, `
		SELECT count(*) FROM ride_offers WHERE ride_id = $1 AND status = 'PENDING'`, rideID)
	if err != nil {
		return err
	}

	if remaining == 0 {
		if _, err := tx.Exec(`
			UPDATE rides SET status = 'NO_DRIVERS_FOUND', updated_at = now()
			WHERE id = $1 AND status = 'REQUESTED'`, rideID); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func GetPendingRidesForDriver(db *sqlx.DB, driverID string) ([]models.Ride, error) {
	var rides []models.Ride
	query := `
		SELECT r.id, r.rider_id, r.driver_id, r.status, r.pickup_lat, r.pickup_lng, r.drop_lat, r.drop_lng,
		       r.requested_at, r.accepted_at, r.updated_at
		FROM rides r
		JOIN ride_offers o ON o.ride_id = r.id
		WHERE o.driver_id = $1 AND o.status = 'PENDING' AND r.status = 'REQUESTED'
		ORDER BY r.requested_at`
	err := db.Select(&rides, query, driverID)
	return rides, err
}
