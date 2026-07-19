package repository

import (
	"TravelBackend/models"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
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
		WHERE distance_km < 5`

	var drivers []models.NearbyDriver
	err := db.Select(&drivers, query, lat, lng)
	return drivers, err
}

func CreateRide(db *sqlx.DB, riderID string, req models.RideRequest) (string, error) {
	var id string
	query := `
		INSERT INTO rides (rider_id, pickup_lat, pickup_lng, drop_lat, drop_lng, status)
		VALUES ($1, $2, $3, $4, $5, 'REQUESTED') RETURNING id`
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

func MarkRideReachAtDest(db *sqlx.DB, rideID, driverID string) (*models.Ride, error) {
	ride, err := GetRideByID(db, rideID)
	if err != nil {
		return nil, err
	}

	if ride.DriverID == nil || *ride.DriverID != driverID {
		return nil, models.ErrRideNotActive
	}

	if ride.Status != models.RideStatusAccepted {
		return nil, models.ErrRideNotActive
	}

	var loc models.DriverLocation
	err = db.Get(&loc, `SELECT driver_id, latitude, longitude, updated_at FROM driver_locations WHERE driver_id = $1`, driverID)
	if err != nil {
		return nil, err
	}

	var distanceMeters float64

	err = db.Get(&distanceMeters, `
		SELECT 6371000 * acos(
			cos(radians($1)) * cos(radians($2)) * cos(radians($3) - radians($4))
			+ sin(radians($1)) * sin(radians($2))
		)`, loc.Latitude, ride.DropLat, loc.Longitude, ride.DropLng)

	if err != nil {
		return nil, err
	}

	if distanceMeters > 50 {
		return nil, models.ErrNotAtDroppingLocation
	}

	_, err = db.Exec(`UPDATE rides SET status = 'REACHED_AT_DESTINATION', updated_at = now() WHERE id = $1`, rideID)
	if err != nil {
		return nil, err
	}

	return GetRideByID(db, rideID)
}

func GetAllDriverRides(db *sqlx.DB, driver_id string) ([]models.Ride, error) {
	var rides []models.Ride

	err := db.Select(&rides, `SELECT id, rider_id, driver_id, status, pickup_lat, pickup_lng, drop_lat, drop_lng, requested_at,accepted_at, updated_at, completed_at FROM rides WHERE driver_id = $1`, driver_id)

	if err != nil {
		return nil, err
	}
	return rides, nil
}

func GetRideStatus(db *sqlx.DB, rideID string) (string, error) {
	//var ride models.Ride

	query := `SELECT status FROM rides WHERE id = $1`
	var status string
	err := db.QueryRow(query, rideID).Scan(&status)
	if err != nil {
		return "", err
	}

	return status, nil
}

func CalculateFare(db *sqlx.DB, ride_ID string, driverID string, status string) (float64, error) {
	tx, err := db.Beginx()
	if err != nil {
		return 0, err
	}
	query := `SELECT pickup_lat, pickup_lng, drop_lat, drop_lng FROM rides WHERE id = $1`
	var locationPoints models.RideRequest
	err = tx.QueryRow(query, ride_ID).Scan(
		&locationPoints.PickupLat, &locationPoints.PickupLng,
		&locationPoints.DropLat, &locationPoints.DropLng,
	)

	if err != nil {
		return 0, err
	}

	lat1, lng1, lat2, lng2 := locationPoints.PickupLat, locationPoints.PickupLng, locationPoints.DropLat, locationPoints.DropLng

	distance := models.CalculateDistance(lat1, lng1, lat2, lng2)
	if distance < 1 {
		return 0, errors.New("distance is too low")
	}

	fare := 40 + (distance * 12)

	defer tx.Rollback()

	if status != models.RideStatusReachedAtDest {
		return 0, errors.New("ride is not completed yet")
	}

	query = `UPDATE rides SET status = $1 WHERE driver_id = $2`
	_, err = tx.Exec(query, models.RideStatusCompleted, driverID)

	query = `UPDATE drivers SET status = $1 WHERE id = $2`
	_, err = tx.Exec(query, models.DriverStatusOnline, driverID)

	query = `UPDATE rides SET fare = $1 WHERE id = $2 AND driver_id = $3 AND status = $4`
	_, err = tx.Exec(query, fare, ride_ID, driverID, status)

	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return fare, nil
}
