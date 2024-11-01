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
		log.Fatalf("error parsing database connection string: %v", err)
	}

	dbPool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("error connecting to the database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
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
		log.Println("Database connection closed successfully")
	}
}

func RunMigrations(migrationDir string) error {
	_, err := DBPool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename TEXT UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("error creating migrations table: %v", err)
	}

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("error reading migration directory %s: %v", migrationDir, err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			var exists bool
			err := DBPool.QueryRow(context.Background(), `
				SELECT EXISTS (
					SELECT 1 FROM migrations WHERE filename = $1
				)
			`, file.Name()).Scan(&exists)
			if err != nil {
				return fmt.Errorf("error checking migration status for %s: %v", file.Name(), err)
			}

			if exists {
				log.Printf("Skipping already applied migration: %s", file.Name())
				continue
			}

			filePath := filepath.Join(migrationDir, file.Name())
			query, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("error reading migration file %s: %v", filePath, err)
			}

			_, err = DBPool.Exec(context.Background(), string(query))
			if err != nil {
				return fmt.Errorf("error applying migration %s: %v", file.Name(), err)
			}

			_, err = DBPool.Exec(context.Background(), `
				INSERT INTO migrations (filename) VALUES ($1)
			`, file.Name())
			if err != nil {
				return fmt.Errorf("error recording applied migration %s: %v", file.Name(), err)
			}

			log.Printf("Migration applied successfully: %s", file.Name())
		}
	}
	return nil
}
