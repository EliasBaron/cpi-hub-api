package domain

import (
	"time"
)

type User struct {
	ID        int
	Name      string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Image     string
}

type UserWithSpaces struct {
	User   *User
	Spaces []*Space
}
