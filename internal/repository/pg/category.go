package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
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
		&categoryByCode.Code); err != nil {
		return err
	}

	if data.ID != "" {
		errMsg := "transaction failed : category already exists"
		return errors.New(errMsg)
	}

	_, err = tx.Exec(ctx,
		insertCategoryQuery,
		data.ID,
		data.Name,
		data.Code,
		data.Description,
		data.IsActive,
		data.CreatedBy,
	)
	if err != nil {
		errMsg := "transaction failed : failed to insert category : " + err.Error()
		return errors.New(errMsg)
	}
	err = tx.Commit(ctx)
	if err != nil {
		errMsg := "transaction failed : failed to commit : " + err.Error()
		return errors.New(errMsg)
	}
	return nil
}
