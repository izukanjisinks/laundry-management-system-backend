package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect(connStr string) error {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return nil
}

func RunMigrations(migrationsDir string) error {
	// Find all *.up.sql files and sort them by filename (001_, 002_, ...)
	pattern := filepath.Join(migrationsDir, "*.up.sql")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return fmt.Errorf("failed to glob migration files: %w", err)
	}
	if len(files) == 0 {
		return fmt.Errorf("no migration files found in %s", migrationsDir)
	}
	sort.Strings(files)

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", filepath.Base(file), err)
		}

		log.Printf("✓ Applied migration: %s", filepath.Base(file))
	}

	log.Println("All migrations completed successfully")
	return nil
}

func MigrationsDir() string {
	// Resolve the migrations directory relative to the working directory
	candidates := []string{
		"migrations",
		"../migrations",
		filepath.Join(executableDir(), "migrations"),
	}
	for _, dir := range candidates {
		if info, err := os.Stat(dir); err == nil && info.IsDir() {
			return dir
		}
	}
	return "migrations"
}

func executableDir() string {
	exe, err := os.Executable()
	if err != nil {
		return "."
	}
	return filepath.Dir(exe)
}

// GetDB returns the active database connection pool.
func GetDB() *sql.DB {
	return db
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}

// NullableString converts an empty string to sql.NullString.
func NullableString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// StringOrEmpty returns the string value of a sql.NullString, or empty string.
func StringOrEmpty(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

