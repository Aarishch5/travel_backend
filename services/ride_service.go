package services

import (
	"sync"

	"github.com/jmoiron/sqlx"

	"TravelBackend/database/repository"
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

func RequestRide(db *sqlx.DB, riderID string, req models.RideRequest) (*models.Ride, error) {

	rideID, err := repository.CreateRide(db, riderID, req)

	if err != nil {
		return nil, err
	}

	drivers, err := repository.FindNearbyDrivers(db, req.PickupLat, req.PickupLng)
	if err != nil {
		return nil, err
	}

	if len(drivers) == 0 {
		if err := repository.MarkRideNoDriversFound(db, rideID); err != nil {
			return nil, err
		}
		return repository.GetRideByID(db, rideID)
	}

	driverIDs := make([]string, len(drivers))
	for i, d := range drivers {
		driverIDs[i] = d.DriverID
	}

	if err := repository.CreateRideOffers(db, rideID, driverIDs); err != nil {
		return nil, err
	}

	return repository.GetRideByID(db, rideID)
}

func AcceptRide(db *sqlx.DB, rideID, driverID string) (*models.Ride, error) {
	mu := lockForRide(rideID)
	mu.Lock()
	defer mu.Unlock()

	if err := repository.AcceptRide(db, rideID, driverID); err != nil {
		return nil, err
	}
	if err := repository.UpdateDriverStatus(db, driverID, models.DriverStatusOnTrip); err != nil {
		return nil, err
	}
	return repository.GetRideByID(db, rideID)
}

func RejectRide(db *sqlx.DB, rideID, driverID string) (*models.Ride, error) {
	mu := lockForRide(rideID)
	mu.Lock()
	defer mu.Unlock()
	if err := repository.RejectRide(db, rideID, driverID); err != nil {
		return nil, err
	}
	return repository.GetRideByID(db, rideID)
}

func ReachAtDest(db *sqlx.DB, rideID, driverID string) (*models.Ride, error) {
	return repository.MarkRideReachAtDest(db, rideID, driverID)
}

func GetAllRides(db *sqlx.DB, driverID string) ([]models.Ride, error) {
	return repository.GetAllDriverRides(db, driverID)
}

func CalculateFare(db *sqlx.DB, rideID string, driverID string) (float64, error) {
	status, err := repository.GetRideStatus(db, rideID)
	if err != nil {
		return 0, err
	}

	fare, err := repository.CalculateFare(db, rideID, driverID, status)
	if err != nil {
		return 0, err
	}

	return fare, nil
}
