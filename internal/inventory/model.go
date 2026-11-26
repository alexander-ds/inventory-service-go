package inventory

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// ErrNotFound permite a la capa HTTP responder 404 sin acoplarse a SQL.
var ErrNotFound = errors.New("recurso no encontrado")

// ErrInvalidInput identifica errores de validación de entrada (HTTP 400).
var ErrInvalidInput = errors.New("entrada inválida")

// Product refleja la tabla products de inventory_db.
type Product struct {
	ID            string     `json:"id"`
	SKU           string     `json:"sku"`
	Barcode       *string    `json:"barcode,omitempty"`
	Name          string     `json:"name"`
	Description   *string    `json:"description,omitempty"`
	CategoryID    *string    `json:"categoryId,omitempty"`
	SupplierID    *string    `json:"supplierId,omitempty"`
	PurchasePrice float64    `json:"purchasePrice"`
	SalePrice     float64    `json:"salePrice"`
	Unit          string     `json:"unit"`
	Active        bool       `json:"active"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// ProductInput son los campos que el cliente puede enviar.
type ProductInput struct {
	SKU           string  `json:"sku"`
	Barcode       *string `json:"barcode,omitempty"`
	Name          string  `json:"name"`
	Description   *string `json:"description,omitempty"`
	CategoryID    *string `json:"categoryId,omitempty"`
	SupplierID    *string `json:"supplierId,omitempty"`
	PurchasePrice float64 `json:"purchasePrice"`
	SalePrice     float64 `json:"salePrice"`
	Unit          string  `json:"unit"`
	Active        bool    `json:"active"`
}

// Validate aplica las reglas mínimas antes de tocar la base de datos.
func (in ProductInput) Validate() error {
	if strings.TrimSpace(in.SKU) == "" {
		return fmt.Errorf("%w: sku es obligatorio", ErrInvalidInput)
	}
	if strings.TrimSpace(in.Name) == "" {
		return fmt.Errorf("%w: name es obligatorio", ErrInvalidInput)
	}
	if strings.TrimSpace(in.Unit) == "" {
		return fmt.Errorf("%w: unit es obligatorio", ErrInvalidInput)
	}
	if in.PurchasePrice < 0 || in.SalePrice < 0 {
		return fmt.Errorf("%w: los precios no pueden ser negativos", ErrInvalidInput)
	}
	return nil
}
