package utils

import (
	"errors"
	"regexp"
)

var (
	number   = regexp.MustCompile(`[0-9]`)
	emailExp = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
)

func ValidateEmail(email string) error {
	if !emailExp.MatchString(email) {
		return errors.New("invalid email address")
	}
	return nil
}

func ValidatePhoneNumber(phone string) error {
	if !number.MatchString(phone) && len(phone) != 10 {
		return errors.New("invalid phone number")
	}
	return nil
}
