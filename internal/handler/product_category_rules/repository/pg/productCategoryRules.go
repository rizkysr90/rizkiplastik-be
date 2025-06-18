package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/repository"
)

var (
	errCategoryNotFound      = errors.New("category not found")
	errPackagingTypeNotFound = errors.New("packaging type not found")
	errRuleAlreadyExists     = errors.New("rule already exists")
	errUniqueViolation       = errors.New("unique violation")
)

type ProductCategoryRules struct {
	db *pgxpool.Pool
}

func NewProductCategoryRules(db *pgxpool.Pool) *ProductCategoryRules {
	return &ProductCategoryRules{
		db: db,
	}
}

const (
	insertProductCategoryRules = `
		INSERT INTO product_category_rules 
		(
			rule_id, 
			category_id,
			packaging_type_id,
			is_default, 
			created_at,
			created_by,
			updated_at,
			updated_by
		)
		VALUES 
		(
			$1, 
			$2, 
			$3, 
			$4, 
			NOW(), 
			$6, 
			NOW(), 
			$8
		)
	`
	checkRuleByCategoryIDAndPackagingTypeID = `
		SELECT 
			rule_id,
		FROM product_category_rules
		WHERE category_id = $1 AND packaging_type_id = $2
	`
	checkCategoryIDSQL = `
		SELECT 
			id
		FROM product_categories
		WHERE id = $1 and is_active = true
	`
	checkPackagingTypeIDSQL = `
		SELECT 
			id
		FROM packaging_types
		WHERE id = $1 and is_active = true
	`
)

func (r *ProductCategoryRules) InsertTransaction(
	ctx context.Context,
	data *repository.ProductCategoryRulesData,
) error {
	tx, err := r.db.BeginTx(
		ctx,
		pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted,
		},
	)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// Validate category id
	if err := r.checkCategoryID(ctx, tx, data.CategoryID); err != nil {
		return err
	}
	// Validate packaging type id
	if err := r.checkPackagingTypeID(ctx, tx, data.PackagingTypeID); err != nil {
		return err
	}
	// Validate existing rule
	if err := r.checkExistingRule(ctx, tx, data); err != nil {
		return err
	}
	// Insert product category rules
	_, err = tx.Exec(
		ctx,
		insertProductCategoryRules,
		data.RuleID,
		data.CategoryID,
		data.PackagingTypeID,
		data.IsDefault,
		data.CreatedBy,
		data.UpdatedBy,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Check for unique violation error code
			if pgErr.Code == constants.ErrCodePostgreUniqueViolation {
				return errUniqueViolation
			}
		}
		return err
	}
	return nil
}
func (r *ProductCategoryRules) checkExistingRule(
	ctx context.Context,
	tx pgx.Tx,
	data *repository.ProductCategoryRulesData,
) error {
	var ruleID string
	err := tx.QueryRow(
		ctx,
		checkRuleByCategoryIDAndPackagingTypeID,
		data.CategoryID,
		data.PackagingTypeID,
	).Scan(&ruleID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if ruleID != "" {
		return errRuleAlreadyExists
	}
	return nil
}
func (r *ProductCategoryRules) checkCategoryID(
	ctx context.Context,
	tx pgx.Tx,
	categoryID string,
) error {
	var productCategoryID string
	err := tx.QueryRow(
		ctx,
		checkCategoryIDSQL,
		categoryID,
	).Scan(&productCategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errCategoryNotFound
		}
		return err
	}
	return nil
}
func (r *ProductCategoryRules) checkPackagingTypeID(
	ctx context.Context,
	tx pgx.Tx,
	inputPackagingTypeID string,
) error {
	var packagingTypeID string
	err := tx.QueryRow(
		ctx,
		checkPackagingTypeIDSQL,
		inputPackagingTypeID,
	).Scan(&packagingTypeID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errPackagingTypeNotFound
		}
		return err
	}
	return nil
}
