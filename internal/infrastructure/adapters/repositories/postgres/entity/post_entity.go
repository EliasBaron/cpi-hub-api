package entity

import "time"

type PostEntity struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Image     *string   `db:"image"`
	CreatedBy int       `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	UpdatedBy int       `db:"updated_by"`
	SpaceID   int       `db:"space_id"`
}
