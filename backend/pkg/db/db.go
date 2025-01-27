package db

import (
  "database/sql"
  "elysium-backend/config"
  "log"
  "os"
  "path/filepath"

  _ "github.com/mattn/go-sqlite3"
)

var DBPool *sql.DB

func InitializeDatabaseConnection() *sql.DB {
  if config.GetLogLevel() == "DEBUG" {
    log.Println("db.InitializeDatabaseConnection -> called")
  }

  dsn := config.GetEnv("DB_NAME", "elysium.db")

  dbPool, err := sql.Open("sqlite3", dsn)
  if err != nil {
    log.Fatalf("db.InitializeDatabaseConnection -> error connecting to the database: %v", err)
  }

  log.Println("db.InitializeDatabaseConnection -> Successfully connected to SQLite database")
  return dbPool
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

  _, err := DBPool.Exec(`
    CREATE TABLE IF NOT EXISTS migrations (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
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
      err := DBPool.QueryRow(`
        SELECT EXISTS (
        SELECT 1 FROM migrations WHERE filename = ?
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

      _, err = DBPool.Exec(string(query))
      if err != nil {
        log.Println("db.RunMigrations -> error applying migration", file.Name(), ":", err)
        return err
      }

      _, err = DBPool.Exec(`
        INSERT INTO migrations (filename) VALUES (?)
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
