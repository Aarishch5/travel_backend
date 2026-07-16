package utils

import (
	"errors"
	"regexp"
)

var (
	upperCase      = regexp.MustCompile(`[A-Z]`)
	lowerCase      = regexp.MustCompile(`[a-z]`)
	numberPassword = regexp.MustCompile(`[0-9]`)
	special        = regexp.MustCompile(`[^a-zA-Z0-9]`)
	number         = regexp.MustCompile(`^[0-9]+$`)
	emailExp       = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func ValidateEmail(email string) error {
	if !emailExp.MatchString(email) {
		return errors.New("invalid email address")
	}
	return nil
}

func ValidatePhoneNumber(phone string) error {
	if !number.MatchString(phone) || len(phone) != 10 {
		return errors.New("invalid phone number")
	}
	return nil
}

func ValidatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if !upperCase.MatchString(password) {
		return errors.New("password must contain at least one uppercase letter")
	}

	if !lowerCase.MatchString(password) {
		return errors.New("password must contain at least one lowercase letter")
	}

	if !numberPassword.MatchString(password) {
		return errors.New("password must contain at least one number")
	}

	if !special.MatchString(password) {
		return errors.New("password must contain at least one special character")
	}

	return nil
}
