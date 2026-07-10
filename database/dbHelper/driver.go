package dbHelper

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"TravelBackend/models"
)

func CreateDriver(db *sqlx.DB, req models.CreateDriverRequest) (string, error) {
	var id string
	query := `
		INSERT INTO drivers (name, email, phone, license_number, vehicle_model, plate_number)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	err := db.Get(&id, query, req.Name, req.Email, req.Phone, req.LicenseNumber,
		req.VehicleModel, req.PlateNumber)
	return id, err
}

func GetDriverByEmailOrPhone(db *sqlx.DB, email, phone string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM drivers WHERE email = $1 OR phone = $2)`
	err := db.Get(&exists, query, email, phone)
	return exists, err
}

func GetDriverByID(db *sqlx.DB, id string) (*models.Driver, error) {
	var d models.Driver
	query := `SELECT id, name, email, phone, license_number, vehicle_model, plate_number, avg_rating, created_at
		FROM drivers WHERE id = $1`
	err := db.Get(&d, query, id)
	if err == sql.ErrNoRows {
		return nil, models.ErrDriverNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func DeleteDriver(db *sqlx.DB, id string) error {
	result, err := db.Exec(`DELETE FROM drivers WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("driver not found")
	}
	return nil
}
