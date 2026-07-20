package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"TravelBackend/models"
)

func CreateDriver(db *sqlx.DB, req models.CreateDriverRequest, passwordHash string) (string, error) {
	var id string
	query := `INSERT INTO drivers (name, email, phone, license_number, vehicle_model, plate_number, password_hash)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	err := db.Get(&id, query, req.Name, req.Email, req.Phone, req.LicenseNumber,
		req.VehicleModel, req.PlateNumber, passwordHash)
	return id, err
}

func GetDriverByEmailOrPhone(db *sqlx.DB, email, phone string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(*)>0 FROM drivers WHERE email = $1 OR phone = $2`
	err := db.Get(&exists, query, email, phone)
	return exists, err
}

func GetDriverByID(db *sqlx.DB, id string) (*models.Driver, error) {
	var d models.Driver
	query := `SELECT id, name, email, phone, license_number, vehicle_model, plate_number,status, created_at
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

func GetDriverByEmail(db *sqlx.DB, email string) (*models.Driver, error) {
	var d models.Driver
	query := `SELECT id, name, email, phone, license_number, vehicle_model, plate_number, status, password_hash, created_at
		FROM drivers WHERE email = $1`
	err := db.Get(&d, query, email)
	if err == sql.ErrNoRows {
		return nil, models.ErrDriverNotFound
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func UpdateDriverStatus(db *sqlx.DB, id, status string) error {

	query := `UPDATE drivers SET status = $1 WHERE id = $2`

	result, err := db.Exec(query, status, id)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return models.ErrDriverNotFound
	}
	return nil
}

func DeleteDriver(db *sqlx.DB, id string) error {

	query := `UPDATE drivers SET archived_at=now() WHERE id = $1`

	result, err := db.Exec(query, id)

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
