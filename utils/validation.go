package utils

import "strings"

func ValidateCreateDriver(name, email, phone, license string) string {
	switch {
	case strings.TrimSpace(name) == "":
		return "name is required"
	case !strings.Contains(email, "@"):
		return "a valid email is required"
	case strings.TrimSpace(phone) == "":
		return "phone is required"
	case strings.TrimSpace(license) == "":
		return "license_number is required"
	}
	return ""
}

func ValidateCreateRider(name, email, phone string) string {
	switch {
	case strings.TrimSpace(name) == "":
		return "name is required"
	case !strings.Contains(email, "@"):
		return "a valid email is required"
	case strings.TrimSpace(phone) == "":
		return "phone is required"
	}
	return ""
}
