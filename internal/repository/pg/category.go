package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
)

// Sentinel errors - define once, use everywhere
var (
	ErrCategoryAlreadyExists = errors.New("category already exists")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrDatabaseOperation     = errors.New("database operation failed")
	ErrTransactionFailed     = errors.New("transaction failed")
)

type Category struct {
	db *pgxpool.Pool
}

func NewCategory(db *pgxpool.Pool) *Category {
	return &Category{db: db}
}

const (
	insertCategoryQuery = `
		INSERT INTO product_categories (
			id, 
			name, 
			code, 
			description, 
			is_active, 
			created_by, 
			updated_by, 
			created_at, 
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`
	findByCodeQuery = `
		SELECT id, name, code
		FROM product_categories
		WHERE code = $1
	`
)

func (c *Category) InsertTransaction(
	ctx context.Context, data *repository.CategoryData) error {
	var description interface{}
	description = nil
	if data.Description != "" {
		description = data.Description
	}
	tx, err := c.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var categoryByCode repository.CategoryData
	row := tx.QueryRow(ctx, findByCodeQuery, data.Code)
	if err := row.Scan(
		&categoryByCode.ID,
		&categoryByCode.Name,
		&categoryByCode.Code); err != nil && err != pgx.ErrNoRows {
		return err
	}
	if categoryByCode.ID != "" {
		return ErrCategoryAlreadyExists
	}

	_, err = tx.Exec(ctx,
		insertCategoryQuery,
		data.ID,
		data.Name,
		data.Code,
		description,
		data.IsActive,
		data.CreatedBy,
		data.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseOperation, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrTransactionFailed, err)
	}
	return nil
}
