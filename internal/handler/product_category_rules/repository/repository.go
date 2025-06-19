package repository

import (
	"context"
	"database/sql"
	"time"
)

type ProductCategoryRulesData struct {
	RuleID          string
	CategoryID      string
	PackagingTypeID string
	IsDefault       bool
	CreatedAt       time.Time
	CreatedBy       string
	UpdatedAt       time.Time
	UpdatedBy       string
	DeletedAt       sql.NullTime
	DeletedBy       sql.NullString
}

type ProductCategoryRules interface {
	InsertTransaction(
		ctx context.Context,
		data *ProductCategoryRulesData,
	) error
	UpdateTransaction(
		ctx context.Context,
		data *ProductCategoryRulesData,
	) error
}
