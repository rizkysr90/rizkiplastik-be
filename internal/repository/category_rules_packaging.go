package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type CategoryPackagingRulesData struct {
	RuleID            string
	ProductCategoryID string
	PackagingTypeID   string
	PackagingTypeCode string
	PackagingTypeName string
	IsDefault         bool
	IsActive          bool
}

type CategoryPackagingRules interface {
	FindByCategoryIDAndRuleID(
		ctx context.Context,
		tx pgx.Tx,
		categoryID string, packagingTypeID []string,
	) ([]CategoryPackagingRulesData, error)
}
