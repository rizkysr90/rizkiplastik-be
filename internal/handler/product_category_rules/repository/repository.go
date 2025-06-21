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
	PackagingCode   string
	PackagingName   string
	IsDefault       bool
	IsActive        bool
	CreatedAt       time.Time
	CreatedBy       string
	UpdatedAt       time.Time
	UpdatedBy       string
	DeletedAt       sql.NullTime
	DeletedBy       sql.NullString
}

type ProductCategoryRulesFilter struct {
	CategoryID string
	Status     string
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
	FindRuleByCategoryID(
		ctx context.Context,
		filter ProductCategoryRulesFilter,
	) ([]ProductCategoryRulesData, error)
	UpdateStatusRule(
		ctx context.Context,
		ruleID string,
		isActive bool,
	) error
}
