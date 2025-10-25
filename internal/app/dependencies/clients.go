package dependencies

import (
	"cpi-hub-api/database/schema"
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
)

type MongoDBConfig struct {
	URI      string
	Database string
	Timeout  time.Duration
}

type PostgreSQLConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	SSLMode  string
}

func NewPostgreSQLClient() (*sql.DB, error) {
	config := PostgreSQLConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "rootroot",
		Database: "cpihub-beta",
		SSLMode:  "disable",
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying connection to PostgreSQL: %w", err)
	}

	if err = schema.EnsureSchema(db); err != nil {
		return nil, fmt.Errorf("error ensuring database schema: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL at %s:%d", config.Host, config.Port)

	return db, nil
}

func GetPostgreSQLDatabase() (*sql.DB, error) {
	client, err := NewPostgreSQLClient()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func ClosePostgreSQLConnection(db *sql.DB) error {
	return db.Close()
}
