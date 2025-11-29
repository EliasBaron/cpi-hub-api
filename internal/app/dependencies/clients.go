package dependencies

import (
	"context"
	"cpi-hub-api/database/schema"
	"fmt"
	"log"
	"os"
	"strconv"
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
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017" // Fallback para desarrollo local
	}

	timeout := 10 * time.Second
	if timeoutStr := os.Getenv("MONGODB_TIMEOUT"); timeoutStr != "" {
		if parsedTimeout, err := time.ParseDuration(timeoutStr); err == nil {
			timeout = parsedTimeout
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error conectando a MongoDB: %w", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error verificando conexión a MongoDB: %w", err)
	}

	log.Printf("Conectado exitosamente a MongoDB")
	return client, nil
}

func GetMongoDatabase() (*mongo.Database, error) {
	client, err := newMongoDBClient()
	if err != nil {
		return nil, err
	}

	databaseName := os.Getenv("MONGODB_DATABASE")
	if databaseName == "" {
		databaseName = "cpihub"
	}

	return client.Database(databaseName), nil
}

func CloseMongoConnection(client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return client.Disconnect(ctx)
}

func NewPostgreSQLClient() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL != "" {
		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			return nil, fmt.Errorf("error opening connection to PostgreSQL: %w", err)
		}

		if err := db.Ping(); err != nil {
			return nil, fmt.Errorf("error verifying connection to PostgreSQL: %w", err)
		}

		if err = schema.EnsureSchema(db); err != nil {
			return nil, fmt.Errorf("error ensuring database schema: %w", err)
		}

		log.Printf("Successfully connected to PostgreSQL using DATABASE_URL")
		return db, nil
	}

	// Fallback: leer variables individuales
	host := getEnv("POSTGRES_HOST", "localhost")
	sslMode := getEnv("POSTGRES_SSLMODE", "")
	if sslMode == "" {
		if host == "localhost" || host == "127.0.0.1" {
			sslMode = "disable"
		} else {
			sslMode = "require"
		}
	}

	config := PostgreSQLConfig{
		Host:     host,
		Port:     getEnvAsInt("POSTGRES_PORT", 5432),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "rootroot"),
		Database: getEnv("POSTGRES_DB", "cpihub"),
		SSLMode:  sslMode,
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

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
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
