package entity

import "time"

type UserEntity struct {
	ID        string    `bson:"_id"`
	Name      string    `bson:"name"`
	LastName  string    `bson:"last_name"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
	Image     string    `bson:"image"`
}

type SimpleUserEntity struct {
	ID       string `bson:"_id"`
	Name     string `bson:"name"`
	LastName string `bson:"last_name"`
	Image    string `bson:"image"`
	Email    string `bson:"email"`
}
