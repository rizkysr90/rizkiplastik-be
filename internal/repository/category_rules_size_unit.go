package repository

import (
	"context"
	"database/sql"
	"time"
)

type ProductSizeUnitRulesData struct {
	RuleID            string
	ProductCategoryID string
	SizeUnitID        string
	SizeUnitCode      string
	SizeUnitName      string
	SizeUnitType      string
	IsDefault         bool
	IsActive          bool
	CreatedAt         time.Time
	CreatedBy         string
	UpdatedAt         time.Time
	UpdatedBy         string
	DeletedAt         sql.NullTime
}
type ProductSizeUnitRules interface {
	InsertTransaction(
		ctx context.Context,
		data *ProductSizeUnitRulesData,
	) error
}
