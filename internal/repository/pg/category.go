package pg

import (
	"context"
	"database/sql"
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
	updateCategoryQuery = `
		UPDATE product_categories
		SET 
			name = $1, 
			description = $2, 
			is_active = $3, 
			updated_by = $4, 
			updated_at = NOW()
		WHERE id = $5
	`
	findListCategoryQuery = `
		SELECT 
		      id, 
			  name, 
			  code, 
			  is_active,
			  created_at,
			  updated_at,
			  COUNT(*) OVER () AS total_count
		FROM product_categories
		WHERE $1 = '' OR name ILIKE '%' || $1 || '%'
		AND $2 = '' OR code ILIKE '%' || $2 || '%'
		AND CASE 
		        WHEN $3 = 'ALL' THEN TRUE 
		        WHEN $3 = 'TRUE' THEN is_active = true
		        WHEN $3 = 'FALSE' THEN is_active = false
		    END
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`
	findByCategoryIDQuery = `
		SELECT 
			id, 
			name, 
			code, 
			description, 
			is_active, 
			created_at, 
			updated_at
		FROM product_categories
		WHERE id = $1
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
		return ErrAlreadyExists
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
		return errors.New("failed to insert category : " + err.Error())
	}
	if err = tx.Commit(ctx); err != nil {
		return errors.New("failed to commit transaction : " + err.Error())
	}
	return nil
}

func (c *Category) Update(ctx context.Context, data *repository.CategoryData) error {
	var description interface{} = nil
	if data.Description != "" {
		description = data.Description
	}
	result, err := c.db.Exec(
		ctx,
		updateCategoryQuery,
		data.Name,
		description,
		data.IsActive,
		data.UpdatedBy,
		data.ID,
	)
	if err != nil {
		return errors.New("failed to update category : " + err.Error())
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
func (c *Category) GetList(
	ctx context.Context, filter *repository.CategoryDataFilter) (
	[]repository.CategoryData, int, error) {
	rows, err := c.db.Query(ctx, findListCategoryQuery,
		filter.CategoryName,
		filter.CategoryCode,
		filter.IsActive,
		filter.PageSize,
		filter.Offset,
	)
	if err != nil {
		return nil, 0, errors.New("failed to get list category : " + err.Error())
	}
	defer rows.Close()
	var categories []repository.CategoryData
	var totalCount int
	for rows.Next() {
		var category repository.CategoryData
		if err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Code,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
			&totalCount,
		); err != nil {
			return nil, 0, errors.New("failed to scan category data : " + err.Error())
		}

		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, errors.New("rows error : " + err.Error())
	}
	return categories, totalCount, nil
}
func (c *Category) GetByID(ctx context.Context, categoryID string) (
	*repository.CategoryData, error) {
	row := c.db.QueryRow(ctx, findByCategoryIDQuery, categoryID)
	var category repository.CategoryData
	var nullDescription sql.NullString
	if err := row.Scan(
		&category.ID,
		&category.Name,
		&category.Code,
		&nullDescription,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	); err != nil {
		return nil, err
	}
	if nullDescription.Valid {
		category.Description = nullDescription.String
	}
	return &category, nil
}
