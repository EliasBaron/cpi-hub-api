package dependencies

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

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

func newMongoDBClient() (*mongo.Client, error) {
	config := MongoDBConfig{
		URI:      "mongodb://localhost:27017",
		Database: "cpihub",
		Timeout:  10 * time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.URI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error connecting to MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error verifying connection to MongoDB: %w", err)
	}

	log.Printf("Successfully connected to MongoDB at %s", config.URI)
	return client, nil
}

func GetMongoDatabase() (*mongo.Database, error) {
	client, err := newMongoDBClient()
	if err != nil {
		return nil, err
	}

	return client.Database("cpihub"), nil
}

func CloseMongoConnection(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Disconnect(ctx)
}

func NewPostgreSQLClient() (*sql.DB, error) {
	config := PostgreSQLConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "rootroot",
		Database: "cpihub",
		SSLMode:  "disable",
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to PostgreSQL: %w", err)
	}

	if err := ensureSchema(db); err != nil {
		return nil, fmt.Errorf("error ensuring schema in PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verifying connection to PostgreSQL: %w", err)
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

func ensureSchema(db *sql.DB) error {
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
		`CREATE TABLE IF NOT EXISTS spaces (
            id SERIAL PRIMARY KEY,
            name TEXT NOT NULL,
            description TEXT NOT NULL,
            created_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            created_at TIMESTAMP NOT NULL DEFAULT now(),
            updated_by INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            updated_at TIMESTAMP NOT NULL DEFAULT now()
        )`,
		`CREATE TABLE IF NOT EXISTS user_spaces (
            user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
            space_id INT NOT NULL REFERENCES spaces(id) ON DELETE CASCADE,
            PRIMARY KEY (user_id, space_id)
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
