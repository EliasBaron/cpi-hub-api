package dto

type LoginResponse struct {
	User  UserDTO `json:"user"`
	Token string  `json:"token"`
}
