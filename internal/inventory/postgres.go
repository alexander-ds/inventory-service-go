package inventory

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const productColumns = `id, sku, barcode, name, description, category_id,
	supplier_id, purchase_price, sale_price, unit, active, created_at, updated_at`

const createProductQuery = `INSERT INTO products
	(sku, barcode, name, description, category_id, supplier_id,
	 purchase_price, sale_price, unit, active)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING ` + productColumns

const updateProductQuery = `UPDATE products SET
	sku=$1, barcode=$2, name=$3, description=$4, category_id=$5,
	supplier_id=$6, purchase_price=$7, sale_price=$8, unit=$9,
	active=$10, updated_at=NOW()
	WHERE id=$11
	RETURNING ` + productColumns

// PostgresRepository implementa Repository contra inventory_db.
type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

type scanner interface {
	Scan(dest ...any) error
}

func scanProduct(row scanner) (*Product, error) {
	var product Product
	err := row.Scan(
		&product.ID,
		&product.SKU,
		&product.Barcode,
		&product.Name,
		&product.Description,
		&product.CategoryID,
		&product.SupplierID,
		&product.PurchasePrice,
		&product.SalePrice,
		&product.Unit,
		&product.Active,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *PostgresRepository) ListProducts(ctx context.Context) ([]Product, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+productColumns+` FROM products ORDER BY sku`)
	if err != nil {
		return nil, fmt.Errorf("listar productos: %w", err)
	}
	defer rows.Close()

	products := make([]Product, 0)
	for rows.Next() {
		product, err := scanProduct(rows)
		if err != nil {
			return nil, fmt.Errorf("leer producto: %w", err)
		}
		products = append(products, *product)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterar productos: %w", err)
	}
	return products, nil
}

func (r *PostgresRepository) GetProduct(ctx context.Context, id string) (*Product, error) {
	query := `SELECT ` + productColumns + ` FROM products WHERE id=$1`
	product, err := scanProduct(r.db.QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("obtener producto: %w", err)
	}
	return product, nil
}

func (r *PostgresRepository) CreateProduct(ctx context.Context, input ProductInput) (*Product, error) {
	product, err := scanProduct(r.db.QueryRowContext(
		ctx,
		createProductQuery,
		input.SKU,
		input.Barcode,
		input.Name,
		input.Description,
		input.CategoryID,
		input.SupplierID,
		input.PurchasePrice,
		input.SalePrice,
		input.Unit,
		input.Active,
	))
	if err != nil {
		return nil, fmt.Errorf("crear producto: %w", err)
	}
	return product, nil
}

func (r *PostgresRepository) UpdateProduct(ctx context.Context, id string, input ProductInput) (*Product, error) {
	product, err := scanProduct(r.db.QueryRowContext(
		ctx,
		updateProductQuery,
		input.SKU,
		input.Barcode,
		input.Name,
		input.Description,
		input.CategoryID,
		input.SupplierID,
		input.PurchasePrice,
		input.SalePrice,
		input.Unit,
		input.Active,
		id,
	))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("actualizar producto: %w", err)
	}
	return product, nil
}

func (r *PostgresRepository) DeleteProduct(ctx context.Context, id string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id=$1`, id)
	if err != nil {
		return fmt.Errorf("eliminar producto: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("verificar eliminación: %w", err)
	}
	if affected == 0 {
		return ErrNotFound
	}
	return nil
}
