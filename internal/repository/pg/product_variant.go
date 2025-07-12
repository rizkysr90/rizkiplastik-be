package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
)

type ProductVariant struct {
	db *pgxpool.Pool
}

func NewProductVariant(db *pgxpool.Pool) *ProductVariant {
	return &ProductVariant{db: db}
}

const (
	findActiveProductVariantByIDQuery = `
		SELECT 
			id
		FROM product_variants
		WHERE id = ANY($1)
		AND is_active = true
		AND deleted_at IS NULL
	`
	insertProductVariantQuery = `
		INSERT INTO product_variants (
			id, 
			product_id, 
			product_name,
			variant_name, 
			full_name, 
			packaging_type_id,
			size_value,
			size_unit_id,
			cost_price,
			selling_price,
			is_active,
			created_by,
			updated_by,
			created_at,
			updated_at
		) VALUES (
		 	$1, 
			$2, 
			$3,
			$4,
			$5,
			$6,
			$7,
			$8,
			$9,
			$10,
			$11, 
			$12, 
			$13, 
			NOW(), 
			NOW()
		)`
	findProductVariantByProductIDQuery = `
		SELECT 
			id
		FROM product_variants
		WHERE product_id = $1
		AND is_active = true
		AND deleted_at IS NULL
	`
	updateVariantForProductTypeSingleQuery = `
		UPDATE product_variants
		SET packaging_type_id = $1, 
		size_value = $2, 
		size_unit_id = $3, 
		cost_price = $4, 
		selling_price = $5, 
		updated_by = $6, 
		product_name = $7,
		full_name = $8,
		updated_at = NOW()
		WHERE id = $9
	`
)

func (p *ProductVariant) FindManyByID(
	ctx context.Context,
	tx pgx.Tx,
	variantIDs []string,
) ([]repository.ProductVariantData, error) {
	rows, err := tx.Query(
		ctx,
		findActiveProductVariantByIDQuery,
		variantIDs,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	variants := []repository.ProductVariantData{}
	for rows.Next() {
		var variant repository.ProductVariantData
		if err := rows.Scan(&variant.ID); err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return variants, nil
}

func (p *ProductVariant) InsertTransaction(
	ctx context.Context,
	tx pgx.Tx,
	data *repository.ProductVariantData,
) error {
	_, err := tx.Exec(
		ctx, insertProductVariantQuery,
		data.ID,
		data.ProductID,
		data.ProductName,
		data.VariantName,
		data.FullName,
		data.PackagingTypeID,
		data.SizeValue,
		data.SizeUnitID,
		data.CostPrice,
		data.SellingPrice,
		data.IsActive,
		data.CreatedBy,
		data.UpdatedBy,
	)
	if err != nil {
		return err
	}
	return nil
}
func (p *ProductVariant) FindByProductID(
	ctx context.Context,
	tx pgx.Tx,
	productID string,
) ([]repository.ProductVariantData, error) {
	rows, err := tx.Query(
		ctx,
		findProductVariantByProductIDQuery,
		productID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	variants := []repository.ProductVariantData{}
	for rows.Next() {
		var variant repository.ProductVariantData
		if err := rows.Scan(&variant.ID); err != nil {
			return nil, err
		}
		variants = append(variants, variant)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return variants, nil
}
func (p *ProductVariant) UpdateVariantForProductTypeSingleTransaction(
	ctx context.Context,
	tx pgx.Tx,
	data *repository.ProductVariantData,
) error {
	_, err := tx.Exec(
		ctx, updateVariantForProductTypeSingleQuery,
		data.PackagingTypeID,
		data.SizeValue,
		data.SizeUnitID,
		data.CostPrice,
		data.SellingPrice,
		data.UpdatedBy,
		data.ProductName,
		data.FullName,
		data.ID,
	)
	if err != nil {
		return err
	}
	return nil
}
