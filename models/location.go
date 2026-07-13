package models

import "time"

type DriverLocation struct {
	DriverID  string    `json:"driver_id" db:"driver_id"`
	Location  string    `json:"location" db:"location"` // Will be read/written as Hex/WKT strings if scanned directly
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
