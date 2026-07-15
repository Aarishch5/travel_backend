package models

import "errors"

var (
	ErrDriverNotFound        = errors.New("driver not found")
	ErrRiderNotFound         = errors.New("rider not found")
	ErrEmailOrPhoneExists    = errors.New("email or phone already registered")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrInvalidDriverStatus   = errors.New("status must be one of ONLINE, OFFLINE, ON_TRIP")
	ErrRideNotFound          = errors.New("ride not found")
	ErrOfferNotFound         = errors.New("no pending offer found for this ride and driver")
	ErrRideAlreadyTaken      = errors.New("ride has already been accepted by another driver")
	ErrRideNotActive         = errors.New("ride is not in an active accepted state")
	ErrNotAtDroppingLocation = errors.New("driver is not close enough to the drop-off location to complete this ride")
)
