package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

// Config son los parámetros de conexión a inventory_db.
type Config struct {
	Host      string
	Port      int
	User      string
	Password  string
	Name      string
	SSLMode   string
}

// Open crea el pool de conexiones y verifica que la BD sea alcanzable.
func Open(cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	database, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("abrir pool de BD: %w", err)
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("conectar a inventory_db: %w", err)
	}

	return database, nil
}
