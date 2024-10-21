package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

var DBPool *pgxpool.Pool

func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}
}

func InitializeDatabaseConnection() *pgxpool.Pool {
	LoadEnv()

	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v\n", err)
	}

	dbPool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v\n", err)
	}

	log.Println("Connected to PostgreSQL database")
	return dbPool
}

func CloseDatabaseConnection() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database connection closed")
	}
}
