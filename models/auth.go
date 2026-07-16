package models

type AuthResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

const (
	RoleDriver = "driver"
	RoleRider  = "rider"
)
