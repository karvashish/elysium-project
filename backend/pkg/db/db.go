package db

import (
	"context"
	"elysium-backend/config"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
)

var DBPool *pgxpool.Pool

func InitializeDatabaseConnection() *pgxpool.Pool {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("db.InitializeDatabaseConnection -> called")
	}

	dsn := constructDSN()

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("db.InitializeDatabaseConnection -> error parsing database connection string: %v", err)
	}

	dbPool, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatalf("db.InitializeDatabaseConnection -> error connecting to the database: %v", err)
	}

	log.Println("db.InitializeDatabaseConnection -> Successfully connected to PostgreSQL database")
	return dbPool
}

func constructDSN() string {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("db.constructDSN -> called")
	}

	dbUser := config.GetEnv("POSTGRES_USER", "")
	dbPassword := config.GetEnv("POSTGRES_PASSWORD", "")
	dbName := config.GetEnv("POSTGRES_DB", "")
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "5432")

	return "postgresql://" + dbUser + ":" + dbPassword + "@" + dbHost + ":" + dbPort + "/" + dbName
}

func CloseDatabaseConnection() {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("db.CloseDatabaseConnection -> called")
	}

	if DBPool != nil {
		DBPool.Close()
		log.Println("db.CloseDatabaseConnection -> Database connection closed successfully")
	}
}

func RunMigrations(migrationDir string) error {
	if config.GetLogLevel() == "DEBUG" {
		log.Println("db.RunMigrations -> called with migrationDir:", migrationDir)
	}

	_, err := DBPool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			filename TEXT UNIQUE,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Println("db.RunMigrations -> error creating migrations table:", err)
		return err
	}

	files, err := os.ReadDir(migrationDir)
	if err != nil {
		log.Println("db.RunMigrations -> error reading migration directory:", err)
		return err
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
				log.Println("db.RunMigrations -> error checking migration status for", file.Name(), ":", err)
				return err
			}

			if exists {
				log.Println("db.RunMigrations -> Skipping already applied migration:", file.Name())
				continue
			}

			filePath := filepath.Join(migrationDir, file.Name())
			query, err := os.ReadFile(filePath)
			if err != nil {
				log.Println("db.RunMigrations -> error reading migration file:", filePath, ":", err)
				return err
			}

			_, err = DBPool.Exec(context.Background(), string(query))
			if err != nil {
				log.Println("db.RunMigrations -> error applying migration", file.Name(), ":", err)
				return err
			}

			_, err = DBPool.Exec(context.Background(), `
				INSERT INTO migrations (filename) VALUES ($1)
			`, file.Name())
			if err != nil {
				log.Println("db.RunMigrations -> error recording applied migration", file.Name(), ":", err)
				return err
			}

			log.Println("db.RunMigrations -> Migration applied successfully:", file.Name())
		}
	}
	return nil
}
