package services

import (
	"strings"

	"github.com/jmoiron/sqlx"

	"TravelBackend/database/repository"
	"TravelBackend/models"
	"TravelBackend/utils"
)

// Creating the new driver

func RegisterDriver(db *sqlx.DB, req models.CreateDriverRequest) (*models.Driver, error) {
	req.Name = strings.TrimSpace(req.Name)

	exists, err := repository.GetDriverByEmailOrPhone(db, req.Email, req.Phone)

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

	id, err := repository.CreateDriver(db, req, hash)
	if err != nil {
		return nil, err
	}

	return repository.GetDriverByID(db, id)
}

//Verifying the registered driver

func LoginDriver(db *sqlx.DB, req models.LoginDriverRequest) (string, *models.Driver, error) {
	driver, err := repository.GetDriverByEmail(db, req.Email)

	if err == models.ErrDriverNotFound {
		return "", nil, models.ErrInvalidCredentials
	}

	if err != nil {
		return "", nil, err
	}

	if !utils.ComparePassword(driver.PasswordHash, req.Password) {
		return "", nil, models.ErrInvalidCredentials
	}

	token, err := utils.GenerateToken(driver.ID, models.RoleDriver)
	if err != nil {
		return "", nil, err
	}

	driver.PasswordHash = ""
	return token, driver, nil
}

// Updating the drivers status

func UpdateDriverStatus(db *sqlx.DB, driverID, status string) error {
	if !models.IsValidDriverStatus(status) {
		return models.ErrInvalidDriverStatus
	}
	return repository.UpdateDriverStatus(db, driverID, status)
}
