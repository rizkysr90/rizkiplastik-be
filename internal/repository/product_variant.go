package repository

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type ProductVariantData struct {
	ID              string
	ProductID       string
	ProductName     string
	VariantName     sql.NullString
	FullName        string
	PackagingTypeID string
	SizeValue       float32
	SizeUnitID      string
	CostPrice       decimal.NullDecimal
	SellingPrice    decimal.Decimal
	IsActive        bool
	CreatedBy       string
	UpdatedBy       string
	RepackRecipe    *RepackRecipeData
}

type ProductVariant interface {
	FindManyByID(
		ctx context.Context,
		tx pgx.Tx,
		variantIDs []string,
	) ([]ProductVariantData, error)
	InsertTransaction(
		ctx context.Context,
		tx pgx.Tx,
		data *ProductVariantData,
	) error
}
