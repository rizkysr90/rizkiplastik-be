package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
)

type Product struct {
	db *pgxpool.Pool
}

func NewProduct(db *pgxpool.Pool) *Product {
	return &Product{db: db}
}

const (
	insertProductQuery = `
		INSERT INTO products (
			id, 
			base_name, 
			category_id, 
			type, 
			created_by, 
			updated_by,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`
	updateProductQuery = `
		UPDATE products
		SET base_name = $1, category_id = $2, updated_by = $3, updated_at = NOW()
		WHERE id = $4
	`
)

func (p *Product) InsertTransaction(
	ctx context.Context,
	tx pgx.Tx,
	data *repository.ProductData,
) error {
	_, err := tx.Exec(
		ctx, insertProductQuery,
		data.ID,
		data.BaseName,
		data.CategoryID,
		data.ProductType,
		data.CreatedBy,
		data.UpdatedBy,
	)
	if err != nil {
		return err
	}
	return nil
}
func (p *Product) UpdateTransaction(
	ctx context.Context,
	tx pgx.Tx,
	data *repository.ProductData,
) error {
	_, err := tx.Exec(
		ctx, updateProductQuery,
		data.BaseName,
		data.CategoryID,
		data.UpdatedBy,
		data.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
