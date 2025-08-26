package domain

import (
	"time"
)

type User struct {
	ID        string
	Name      string
	LastName  string
	Email     string
	Password  string
	CreatedAt time.Time
	UpdatedAt time.Time
	Image     string
}
