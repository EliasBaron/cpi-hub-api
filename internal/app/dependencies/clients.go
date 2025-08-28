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
	connStr := "user=postgres password=rootroot dbname=cpihub sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error conectando a PostgreSQL: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error verificando conexión a PostgreSQL: %w", err)
	}

	log.Printf("Conectado exitosamente a PostgreSQL")
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
