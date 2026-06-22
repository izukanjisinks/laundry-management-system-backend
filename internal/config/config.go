package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	DBSSLMode    string
	DatabaseURL  string
	JWTSecret    string
	JWTExpiry    int
	Port         string
	Environment  string
	AdminEmail   string
	AdminPassword string
}

var config *Config

func Load() *Config {
	// Load .env file if it exists
	godotenv.Load()

	jwtExpiry := 900 // default 15 minutes
	if exp := os.Getenv("JWT_EXPIRY"); exp != "" {
		if parsed, err := strconv.Atoi(exp); err == nil {
			jwtExpiry = parsed
		}
	}

	cfg := &Config{
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "laundry_user"),
		DBPassword:    getEnv("DB_PASSWORD", "laundry_password"),
		DBName:        getEnv("DB_NAME", "laundry_db"),
		DBSSLMode:     getEnv("DB_SSLMODE", "disable"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-this"),
		JWTExpiry:     jwtExpiry,
		Port:          getEnv("PORT", "8080"),
		Environment:   getEnv("ENVIRONMENT", "development"),
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@laundry.local"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "defaultpassword123"),
	}

	// Build DatabaseURL
	cfg.DatabaseURL = fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)

	config = cfg
	return cfg
}

func Get() *Config {
	if config == nil {
		return Load()
	}
	return config
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
