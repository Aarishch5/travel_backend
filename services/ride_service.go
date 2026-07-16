package services

import (
	"sync"

	"github.com/jmoiron/sqlx"

	"TravelBackend/database/dbHelper"
	"TravelBackend/models"
)

var (
	rideLocks   = make(map[string]*sync.Mutex)
	rideLocksMu sync.Mutex
)

func lockForRide(rideID string) *sync.Mutex {
	rideLocksMu.Lock()
	defer rideLocksMu.Unlock()

	if _, ok := rideLocks[rideID]; !ok {
		rideLocks[rideID] = &sync.Mutex{}
	}
	return rideLocks[rideID]
}

func RequestRide(db *sqlx.DB, riderID string, req models.RequestRideRequest) (*models.Ride, error) {
	rideID, err := dbHelper.CreateRide(db, riderID, req)
	if err != nil {
		return nil, err
	}

	drivers, err := dbHelper.FindNearbyDrivers(db, req.PickupLat, req.PickupLng)
	if err != nil {
		return nil, err
	}

	if len(drivers) == 0 {
		if err := dbHelper.MarkRideNoDriversFound(db, rideID); err != nil {
			return nil, err
		}
		return dbHelper.GetRideByID(db, rideID)
	}

	driverIDs := make([]string, len(drivers))
	for i, d := range drivers {
		driverIDs[i] = d.DriverID
	}

	if err := dbHelper.CreateRideOffers(db, rideID, driverIDs); err != nil {
		return nil, err
	}

	return dbHelper.GetRideByID(db, rideID)
}

func AcceptRide(db *sqlx.DB, rideID, driverID string) (*models.Ride, error) {
	mu := lockForRide(rideID)
	mu.Lock()
	defer mu.Unlock()

	if err := dbHelper.AcceptRide(db, rideID, driverID); err != nil {
		return nil, err
	}
	if err := dbHelper.UpdateDriverStatus(db, driverID, models.DriverStatusOnTrip); err != nil {
		return nil, err
	}
	return dbHelper.GetRideByID(db, rideID)
}

func RejectRide(db *sqlx.DB, rideID, driverID string) (*models.Ride, error) {
	mu := lockForRide(rideID)
	mu.Lock()
	defer mu.Unlock()
	if err := dbHelper.RejectRide(db, rideID, driverID); err != nil {
		return nil, err
	}
	return dbHelper.GetRideByID(db, rideID)
}

func CompleteRide(db *sqlx.DB, rideID, driverID string) (*models.Ride, error) {
	return dbHelper.MarkRideCompleted(db, rideID, driverID)
}

func GetAllRides(db *sqlx.DB, driverID string) ([]*models.Ride, error) {
	return dbHelper.GetAllDriverRides(db, driverID)
}
