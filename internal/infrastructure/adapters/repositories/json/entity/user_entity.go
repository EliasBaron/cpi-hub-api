package entity

import "time"

type UserEntity struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Image     string    `json:"image"`
}
