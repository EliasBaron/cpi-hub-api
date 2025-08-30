package entity

import "time"

type UserEntity struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Image     string    `db:"image"`
}
