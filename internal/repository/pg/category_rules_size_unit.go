package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
)

var (
	ErrRuleSizeUnitAlreadyExists = errors.New("rule size unit already exists")
	ErrRuleSizeUnitNotFound      = errors.New("rule size unit not found")
	ErrUniqueViolation           = errors.New("unique violation")
)

const (
	insertProductSizeUnitRules = `
		INSERT INTO product_categories_size_unit_rules (
			rule_id,
			category_id,
			size_unit_id,
			is_default,
			created_at,
			created_by,
			updated_at,
			updated_by
		)
		VALUES ($1, $2, $3, $4, NOW(), $5, NOW(), $6)
	`
	checkSizeUnitIDSQL = `
		SELECT 
			id
		FROM size_units
		WHERE id = $1 and is_active = true
	`
	checkRuleByCategoryIDAndSizeUnitIDSQL = `
		SELECT 
			rule_id
		FROM product_categories_size_unit_rules
		WHERE category_id = $1 AND size_unit_id = $2
	`
	updateRuleSizeUnitSQL = `
		UPDATE product_categories_size_unit_rules
		SET
			size_unit_id = $2,
			is_default = $3,
			updated_at = NOW(),
			updated_by = $4
		WHERE rule_id = $1
		AND is_active = true
		AND category_id = $5
	`
	findSizeUnitRulesByCategoryIDSQL = `
		SELECT
			a.rule_id,
			a.category_id,
			a.size_unit_id,
			a.is_default,
			a.is_active,
			b.code,
			b.name,
			b.unit_type
		FROM product_categories_size_unit_rules a
		JOIN size_units b
			ON b.id = a.size_unit_id
		WHERE a.category_id = $1
		AND (
			CASE
				WHEN $2 = 'TRUE' THEN
					a.is_active = true
				WHEN $2 = 'FALSE' THEN
					a.is_active = false
				ELSE
				  	true
			END
		)
	`
	updateStatusRuleSizeUnitSQL = `
		UPDATE product_categories_size_unit_rules
		SET
			is_active = $2,
			updated_at = NOW(),
			updated_by = $3
		WHERE rule_id = $1
	`
	findSizeUnitRuleByCategoryIDAndRuleIDSQL = `
		SELECT
			p.category_id,
			p.size_unit_id,
			s.code,
			pc.code
		FROM product_categories_size_unit_rules p
		JOIN size_units s
			ON p.size_unit_id = s.id
		JOIN product_categories pc
			ON pc.id = p.category_id
		WHERE 
			p.category_id = $1 AND 
			p.size_unit_id = ANY($2::uuid[]) AND 
			p.is_active = true AND
			s.is_active = true AND
			pc.is_active = true
	`
)

type ProductSizeUnitRules struct {
	db              *pgxpool.Pool
	productCategory *ProductCategory
	sizeUnit        *SizeUnit
}

func NewProductSizeUnitRules(db *pgxpool.Pool, productCategory *ProductCategory, sizeUnit *SizeUnit) *ProductSizeUnitRules {
	return &ProductSizeUnitRules{db: db, productCategory: productCategory, sizeUnit: sizeUnit}
}

func (pg *ProductSizeUnitRules) InsertTransaction(
	ctx context.Context,
	data *repository.ProductSizeUnitRulesData,
) error {
	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// Validate category id
	if err := pg.productCategory.CheckCategoryID(ctx, tx, data.ProductCategoryID); err != nil {
		return err
	}
	// Validate size unit id
	if err := pg.sizeUnit.checkSizeUnitID(ctx, tx, data.SizeUnitID); err != nil {
		return err
	}
	// Validate existing rule
	if err := pg.checkExistingRule(ctx, tx, data); err != nil {
		return err
	}
	// Insert data
	_, err = tx.Exec(
		ctx,
		insertProductSizeUnitRules,
		data.RuleID,
		data.ProductCategoryID,
		data.SizeUnitID,
		data.IsDefault,
		data.CreatedBy,
		data.UpdatedBy,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == constants.ErrCodePostgreUniqueViolation {
				return ErrRuleSizeUnitAlreadyExists
			}
		}
		return err
	}
	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
func (pg *ProductSizeUnitRules) UpdateTransaction(
	ctx context.Context,
	data *repository.ProductSizeUnitRulesData,
) error {
	tx, err := pg.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	// Validate category id
	if err := pg.productCategory.CheckCategoryID(ctx, tx, data.ProductCategoryID); err != nil {
		return err
	}
	// Validate size unit id
	if err := pg.sizeUnit.checkSizeUnitID(ctx, tx, data.SizeUnitID); err != nil {
		return err
	}
	// Validate existing rule
	if err := pg.checkExistingRule(ctx, tx, data); err != nil {
		return err
	}
	// Update data
	_, err = tx.Exec(
		ctx,
		updateRuleSizeUnitSQL,
		data.RuleID,
		data.SizeUnitID,
		data.IsDefault,
		data.UpdatedBy,
		data.ProductCategoryID,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == constants.ErrCodePostgreUniqueViolation {
				return ErrUniqueViolation
			}
		}
		return err
	}
	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
func (r *ProductSizeUnitRules) checkExistingRule(
	ctx context.Context,
	tx pgx.Tx,
	data *repository.ProductSizeUnitRulesData,
) error {
	var ruleID string
	err := tx.QueryRow(
		ctx,
		checkRuleByCategoryIDAndSizeUnitIDSQL,
		data.ProductCategoryID,
		data.SizeUnitID,
	).Scan(&ruleID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if ruleID != "" && data.RuleID != ruleID {
		return ErrRuleSizeUnitAlreadyExists
	}
	return nil
}
func (pg *ProductSizeUnitRules) FindSizeUnitRulesByCategoryID(
	ctx context.Context,
	filter repository.ProductSizeUnitRulesFilter,
) ([]repository.ProductSizeUnitRulesData, error) {
	rows, err := pg.db.Query(
		ctx,
		findSizeUnitRulesByCategoryIDSQL,
		filter.CategoryID,
		filter.Status,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []repository.ProductSizeUnitRulesData
	for rows.Next() {
		var rule repository.ProductSizeUnitRulesData
		if err := rows.Scan(
			&rule.RuleID,
			&rule.ProductCategoryID,
			&rule.SizeUnitID,
			&rule.IsDefault,
			&rule.IsActive,
			&rule.SizeUnitCode,
			&rule.SizeUnitName,
			&rule.SizeUnitType,
		); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	return rules, nil
}
func (pg *ProductSizeUnitRules) UpdateStatusRule(
	ctx context.Context,
	ruleID string,
	isActive bool,
	userID string,
) error {
	row, err := pg.db.Exec(
		ctx,
		updateStatusRuleSizeUnitSQL,
		ruleID,
		isActive,
		userID,
	)
	if err != nil {
		return err
	}
	if row.RowsAffected() == 0 {
		return ErrRuleSizeUnitNotFound
	}
	return nil
}
func (pg *ProductSizeUnitRules) FindByCategoryIDAndSizeUnitID(
	ctx context.Context,
	tx pgx.Tx,
	categoryID string, sizeUnitID []string,
) ([]repository.ProductSizeUnitRulesData, error) {
	var rules []repository.ProductSizeUnitRulesData
	rows, err := tx.Query(
		ctx,
		findSizeUnitRuleByCategoryIDAndRuleIDSQL,
		categoryID,
		sizeUnitID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rule repository.ProductSizeUnitRulesData
		err := rows.Scan(
			&rule.ProductCategoryID,
			&rule.SizeUnitID,
			&rule.SizeUnitCode,
			&rule.ProductCategoryCode,
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return rules, nil
}
