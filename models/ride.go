package models

import "time"

const (
	RideStatusRequested = "REQUESTED"
	RideStatusAccepted  = "ACCEPTED"
	RideStatusCancelled = "CANCELLED"
	RideStatusCompleted = "COMPLETED"
	RideStatusNoDrivers = "NO_DRIVERS_FOUND"
)

const (
	OfferStatusPending  = "PENDING"
	OfferStatusAccepted = "ACCEPTED"
	OfferStatusRejected = "REJECTED"
	OfferStatusExpired  = "EXPIRED"
)

type Ride struct {
	ID          string     `json:"id" db:"id"`
	RiderID     string     `json:"rider_id" db:"rider_id"`
	DriverID    *string    `json:"driver_id,omitempty" db:"driver_id"`
	Status      string     `json:"status" db:"status"`
	PickupLat   float64    `json:"pickup_lat" db:"pickup_lat"`
	PickupLng   float64    `json:"pickup_lng" db:"pickup_lng"`
	DropLat     float64    `json:"drop_lat" db:"drop_lat"`
	DropLng     float64    `json:"drop_lng" db:"drop_lng"`
	RequestedAt time.Time  `json:"requested_at" db:"requested_at"`
	AcceptedAt  *time.Time `json:"accepted_at,omitempty" db:"accepted_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

type RequestRideRequest struct {
	PickupLat float64 `json:"pickup_lat"`
	PickupLng float64 `json:"pickup_lng"`
	DropLat   float64 `json:"drop_lat"`
	DropLng   float64 `json:"drop_lng"`
}

type NearbyDriver struct {
	DriverID   string  `json:"driver_id" db:"driver_id"`
	DistanceKM float64 `json:"distance_km" db:"distance_km"`
}

type RideOffer struct {
	ID        string    `json:"id" db:"id"`
	RideID    string    `json:"ride_id" db:"ride_id"`
	DriverID  string    `json:"driver_id" db:"driver_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
