package models

import "time"

type DriverLocation struct {
	DriverID  string    `json:"driver_id" db:"driver_id"`
	Latitude  float64   `json:"latitude" db:"latitude"`
	Longitude float64   `json:"longitude" db:"longitude"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
