package dbHelper

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"TravelBackend/models"
)

func CreateRider(db *sqlx.DB, req models.CreateRiderRequest, passwordHash string) (string, error) {
	var id string
	query := `
		INSERT INTO riders (name, email, phone, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	err := db.Get(&id, query, req.Name, req.Email, req.Phone, passwordHash)
	return id, err
}

func GetRiderByEmailOrPhone(db *sqlx.DB, email, phone string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM riders WHERE email = $1 OR phone = $2)`
	err := db.Get(&exists, query, email, phone)
	return exists, err
}

func GetRiderByID(db *sqlx.DB, id string) (*models.Rider, error) {
	var r models.Rider
	query := `SELECT id, name, email, phone, created_at FROM riders WHERE id = $1`
	err := db.Get(&r, query, id)
	if err == sql.ErrNoRows {
		return nil, models.ErrRiderNotFound
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func GetRiderByEmail(db *sqlx.DB, email string) (*models.Rider, error) {
	var r models.Rider
	query := `SELECT id, name, email, phone, password_hash, created_at FROM riders WHERE email = $1`
	err := db.Get(&r, query, email)
	if err == sql.ErrNoRows {
		return nil, models.ErrRiderNotFound
	}
	if err != nil {
		return nil, err
	}
	return &r, nil
}

func DeleteRider(db *sqlx.DB, id string) error {
	result, err := db.Exec(`UPDATE riders SET archived_at=now() WHERE id = $1`, id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("rider not found")
	}
	return nil
}
