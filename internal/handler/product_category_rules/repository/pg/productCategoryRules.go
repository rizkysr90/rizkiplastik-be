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
	ErrCategoryNotFound      = errors.New("category not found")
	ErrPackagingTypeNotFound = errors.New("packaging type not found")
	ErrRuleAlreadyExists     = errors.New("rule already exists")
	ErrUniqueViolation       = errors.New("unique violation")
	ErrRuleNotFound          = errors.New("rule not found")
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
		INSERT INTO product_categories_packaging_rules 
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
			$5, 
			NOW(), 
			$6
		)
	`
	checkRuleByCategoryIDAndPackagingTypeID = `
		SELECT 
			rule_id
		FROM product_categories_packaging_rules
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
	updateRuleSQL = `
		UPDATE 
			product_categories_packaging_rules
		SET
			is_default = $2,
			packaging_type_id = $3,
			updated_at = NOW(),
			updated_by = $4
		WHERE rule_id = $1 
		AND is_active = true
		AND category_id = $5
	`
	findRuleByCategoryID = `
		SELECT 
			pcpr.rule_id,
			pcpr.category_id,
			pcpr.packaging_type_id,
			pcpr.is_default,
			pcpr.is_active,
			pt.code,
			pt.name
		FROM product_categories_packaging_rules pcpr
		JOIN packaging_types pt
			ON pt.id = pcpr.packaging_type_id
		WHERE pcpr.category_id = $1
		AND (
			CASE
				WHEN $2 = 'TRUE' THEN
					pcpr.is_active = true
				WHEN $2 = 'FALSE' THEN
					pcpr.is_active = false
				ELSE
					pcpr.is_active = true
			END
		)
	`
	updateStatusRuleSQL = `
		UPDATE 
			product_categories_packaging_rules
		SET
			is_active = $2
		WHERE rule_id = $1
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
				return ErrUniqueViolation
			}
		}
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
func (r *ProductCategoryRules) UpdateTransaction(
	ctx context.Context,
	data *repository.ProductCategoryRulesData,
) error {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
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
	// Update product category rules
	result, err := tx.Exec(
		ctx,
		updateRuleSQL,
		data.RuleID,
		data.IsDefault,
		data.PackagingTypeID,
		data.UpdatedBy,
		data.CategoryID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// Check for unique violation error code
			if pgErr.Code == constants.ErrCodePostgreUniqueViolation {
				return ErrUniqueViolation
			}
		}
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrRuleNotFound
	}
	if err := tx.Commit(ctx); err != nil {
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
	if ruleID != "" && data.RuleID != ruleID {
		return ErrRuleAlreadyExists
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
			return ErrCategoryNotFound
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
			return ErrPackagingTypeNotFound
		}
		return err
	}
	return nil
}
func (r *ProductCategoryRules) FindRuleByCategoryID(
	ctx context.Context,
	filter repository.ProductCategoryRulesFilter) (
	[]repository.ProductCategoryRulesData, error) {
	rows, err := r.db.Query(
		ctx,
		findRuleByCategoryID,
		filter.CategoryID,
		filter.Status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []repository.ProductCategoryRulesData
	for rows.Next() {
		var rule repository.ProductCategoryRulesData
		if err := rows.Scan(
			&rule.RuleID,
			&rule.CategoryID,
			&rule.PackagingTypeID,
			&rule.IsDefault,
			&rule.IsActive,
			&rule.PackagingCode,
			&rule.PackagingName,
		); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *ProductCategoryRules) UpdateStatusRule(
	ctx context.Context,
	ruleID string,
	isActive bool,
) error {
	row, err := r.db.Exec(
		ctx,
		updateStatusRuleSQL,
		ruleID,
		isActive,
	)
	if err != nil {
		return err
	}
	if row.RowsAffected() == 0 {
		return ErrRuleNotFound
	}
	return nil
}
