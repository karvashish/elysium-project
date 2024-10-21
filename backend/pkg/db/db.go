package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DBPool *pgxpool.Pool

func InitializeDatabaseConnection() *pgxpool.Pool {
	dsn := constructDSN()

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

func constructDSN() string {
	dbUser := getEnv("POSTGRES_USER", "")
	dbPassword := getEnv("POSTGRES_PASSWORD", "")
	dbName := getEnv("POSTGRES_DB", "")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")

	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func CloseDatabaseConnection() {
	if DBPool != nil {
		DBPool.Close()
		log.Println("Database connection closed")
	}
}

func RunMigrations(migrationDir string) error {
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			filePath := filepath.Join(migrationDir, file.Name())
			query, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			_, err = DBPool.Exec(context.Background(), string(query))
			if err != nil {
				return err
			}

			log.Printf("Migration applied: %s\n", file.Name())
		}
	}
	return nil
}
