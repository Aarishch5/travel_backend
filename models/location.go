package models

import "math"

type DriverLocation struct {
	ID        string  `json:"driver_id" db:"driver_id"`
	Latitude  float64 `json:"latitude" db:"latitude"`
	Longitude float64 `json:"longitude" db:"longitude"`
	UpdatedAt string  `json:"updated_at" db:"updated_at"`
}

func degreeToRad(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func CalculateDistance(lat1, lng1, lat2, lng2 float64) float64 {

	ans := 6371.0 * math.Acos(
		math.Cos(degreeToRad(lat1))*
			math.Cos(degreeToRad(lat2))*
			math.Cos(degreeToRad(lng2)-degreeToRad(lng1))+
			math.Sin(degreeToRad(lat1))*
				math.Sin(degreeToRad(lat2)),
	)

	return ans
}
