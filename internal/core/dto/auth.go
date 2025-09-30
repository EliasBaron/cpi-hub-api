package dto

type AuthResponse struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}
