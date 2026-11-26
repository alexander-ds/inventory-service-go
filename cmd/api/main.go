package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexander-ds/EjercicioIA-inventario/inventory-service-go/internal/api"
	"github.com/alexander-ds/EjercicioIA-inventario/inventory-service-go/internal/auth"
	"github.com/alexander-ds/EjercicioIA-inventario/inventory-service-go/internal/config"
	"github.com/alexander-ds/EjercicioIA-inventario/inventory-service-go/internal/db"
	"github.com/alexander-ds/EjercicioIA-inventario/inventory-service-go/internal/inventory"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("configuración inválida: %v", err)
	}

	database, err := db.Open(db.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		Name:     cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		log.Fatalf("no se pudo conectar a inventory_db: %v", err)
	}
	defer database.Close()

	repo := inventory.NewPostgresRepository(database)
	service := inventory.NewService(repo)
	validator := auth.NewValidator(cfg.JWTSecret)
	handler := api.NewServer(service, validator, "inventory-service-go").Handler()

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.Port),
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("inventory-service-go escuchando en %s", server.Addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("servidor HTTP: %v", err)
	}
}
