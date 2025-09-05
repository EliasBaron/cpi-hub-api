package schema

import (
	"database/sql"
	"fmt"
)

func EnsureSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            last_name TEXT NOT NULL,
            email TEXT UNIQUE NOT NULL,
            password TEXT NOT NULL,
            created_at TIMESTAMP NOT NULL DEFAULT now(),
            updated_at TIMESTAMP NOT NULL DEFAULT now(),
            image TEXT
        )`,

		`CREATE TABLE IF NOT EXISTS user_spaces (
            user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            space_id INT NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
            PRIMARY KEY (user_id, space_id)
        )`,
		`CREATE TABLE IF NOT EXISTS spaces (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT NOT NULL,
            created_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            created_at TIMESTAMP NOT NULL DEFAULT now(),
            updated_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            updated_at TIMESTAMP NOT NULL DEFAULT now()
        )`,
		`CREATE TABLE IF NOT EXISTS posts (
            id SERIAL PRIMARY KEY,
            title TEXT NOT NULL,
            content TEXT NOT NULL,
            created_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            created_at TIMESTAMP NOT NULL DEFAULT now(),
            updated_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            updated_at TIMESTAMP NOT NULL DEFAULT now(),
            space_id INT NOT NULL REFERENCES spaces(id) ON DELETE CASCADE
        )`,
		`CREATE TABLE IF NOT EXISTS comments (
			id SERIAL PRIMARY KEY,
			post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
			content TEXT NOT NULL,
    		created_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    		created_at TIMESTAMP NOT NULL DEFAULT now()
		)`,
	}

	for _, stmt := range stmts {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("error creating table: %w", err)
		}
	}
	return nil
}
