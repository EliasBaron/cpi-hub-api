package dependencies

import (
	"context"
	"cpi-hub-api/database/schema"
	"fmt"
	"log"
	"time"

	"database/sql"

	_ "github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		return nil, fmt.Errorf("error conectando a MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error verificando conexión a MongoDB: %w", err)
	}

	log.Printf("Conectado exitosamente a MongoDB en %s", config.URI)
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
