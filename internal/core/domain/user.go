package domain

import (
	"time"
)

const (
	AddUserToSpace      = "add"
	RemoveUserFromSpace = "remove"
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
