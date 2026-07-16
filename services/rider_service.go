package services

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"TravelBackend/database/dbHelper"
	"TravelBackend/models"
	"TravelBackend/utils"
)

// Creating the new rider

func RegisterRider(db *sqlx.DB, req models.CreateRiderRequest) (*models.Rider, error) {
	req.Name = strings.TrimSpace(req.Name)

	exists, err := dbHelper.GetRiderByEmailOrPhone(db, req.Email, req.Phone)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, models.ErrEmailOrPhoneExists
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	id, err := dbHelper.CreateRider(db, req, hash)
	if err != nil {
		return nil, err
	}

	return dbHelper.GetRiderByID(db, id)
}

// validating the registered driver

func LoginRider(db *sqlx.DB, req models.LoginRiderRequest) (string, *models.Rider, error) {
	rider, err := dbHelper.GetRiderByEmail(db, req.Email)
	if err == models.ErrRiderNotFound {
		return "", nil, models.ErrInvalidCredentials
	}
	if err != nil {
		return "", nil, err
	}

	if !utils.ComparePassword(rider.PasswordHash, req.Password) {
		return "", nil, models.ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(rider.ID, models.RoleRider)
	if err != nil {
		return "", nil, err
	}

	rider.PasswordHash = ""
	return token, rider, nil
}
