package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
)

type CategoryPackagingRules struct {
	db *pgxpool.Pool
}

func NewCategoryPackagingRules(db *pgxpool.Pool) *CategoryPackagingRules {
	return &CategoryPackagingRules{
		db: db,
	}
}

const (
	findPackagingRuleByCategoryIDAndRuleIDSQL = `
		SELECT
			p.category_id,
			p.packaging_type_id,
			pt.code
		FROM product_categories_packaging_rules p 
		JOIN packaging_types pt
			ON pt.id = p.packaging_type_id
		JOIN product_categories pc
			ON pc.id = p.category_id
		WHERE 
			p.category_id = $1 AND 
			p.packaging_type_id = ANY($2::uuid[]) AND 
			p.is_active = true AND
			pt.is_active = true AND
			pc.is_active = true
	`
)

func (pg *CategoryPackagingRules) FindByCategoryIDAndRuleID(
	ctx context.Context,
	tx pgx.Tx,
	categoryID string, packagingTypeID []string,
) ([]repository.CategoryPackagingRulesData, error) {
	var rules []repository.CategoryPackagingRulesData
	rows, err := tx.Query(
		ctx,
		findPackagingRuleByCategoryIDAndRuleIDSQL,
		categoryID,
		packagingTypeID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var rule repository.CategoryPackagingRulesData
		err := rows.Scan(
			&rule.ProductCategoryID,
			&rule.PackagingTypeID,
			&rule.PackagingTypeCode,
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
