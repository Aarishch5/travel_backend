package models

import "time"

type Rider struct {
	ID           string    `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	Phone        string    `json:"phone" db:"phone"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	PasswordHash string    `json:"-" db:"password_hash"`
}

type CreateRiderRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type LoginRiderRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
