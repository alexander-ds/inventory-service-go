package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config contiene los valores dependientes del entorno. Nunca se loguea.
type Config struct {
	Port      int
	DBHost    string
	DBPort    int
	DBUser    string
	DBPassword string
	DBName    string
	DBSSLMode string
	JWTSecret string
}

// Load lee la configuración desde variables de entorno.
func Load() (Config, error) {
	port, err := strconv.Atoi(env("PORT", "8080"))
	if err != nil {
		return Config{}, fmt.Errorf("PORT inválido: %w", err)
	}

	dbPort, err := strconv.Atoi(env("DB_PORT", "5432"))
	if err != nil {
		return Config{}, fmt.Errorf("DB_PORT inválido: %w", err)
	}

	cfg := Config{
		Port:       port,
		DBHost:     env("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     env("DB_USER", "inventory_user"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     env("DB_NAME", "inventory_db"),
		DBSSLMode:  env("DB_SSLMODE", "disable"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET es obligatorio")
	}

	return cfg, nil
}

func env(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
