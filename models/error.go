package models

import "errors"

var (
	ErrDriverNotFound = errors.New("driver not found")
	ErrRiderNotFound  = errors.New("rider not found")
)
