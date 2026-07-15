package models

type DriverLocation struct {
	ID        string  `json:"driver_id" db:"driver_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
}
