package models

import "time"

type Driver struct {
	ID            string    `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Email         string    `json:"email" db:"email"`
	Phone         string    `json:"phone" db:"phone"`
	LicenseNumber string    `json:"license_number" db:"license_number"`
	VehicleModel  string    `json:"vehicle_model" db:"vehicle_model"`
	PlateNumber   string    `json:"plate_number" db:"plate_number"`
	AvgRating     float64   `json:"avg_rating" db:"avg_rating"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type CreateDriverRequest struct {
	Name          string `json:"name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	LicenseNumber string `json:"license_number"`
	VehicleModel  string `json:"vehicle_model"`
	PlateNumber   string `json:"plate_number"`
}
