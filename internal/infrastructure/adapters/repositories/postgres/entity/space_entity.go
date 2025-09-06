package entity

import "time"

type SpaceEntity struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedBy   int       `db:"created_by"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedBy   int       `db:"updated_by"`
	UpdatedAt   time.Time `db:"updated_at"`
}
