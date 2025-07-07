package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/shopspring/decimal"
)

type RepackRecipeData struct {
	ID                string
	ParentVariantID   string
	ChildVariantID    string
	QuantityRatio     float32
	RepackCostPerUnit decimal.Decimal
	RepackTimeMinutes int
	CreatedBy         string
	UpdatedBy         string
}

type RepackRecipe interface {
	InsertTransaction(
		ctx context.Context,
		tx pgx.Tx,
		data *RepackRecipeData,
	) error
}
