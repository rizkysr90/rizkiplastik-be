package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
)

type RepackRecipe struct {
	db *pgxpool.Pool
}

func NewRepackRecipe(db *pgxpool.Pool) *RepackRecipe {
	return &RepackRecipe{db: db}
}

const (
	insertRepackRecipeQuery = `
		INSERT INTO product_repack_recipes (
			id,
			parent_variant_id,
			child_variant_id,
			quantity_ratio,
			repack_cost_per_unit,
			repack_time_minutes,
			created_by,
			updated_by,
			created_at,
			updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
	`
)

func (r *RepackRecipe) InsertTransaction(
	ctx context.Context,
	tx pgx.Tx,
	data *repository.RepackRecipeData,
) error {
	_, err := tx.Exec(
		ctx, insertRepackRecipeQuery,
		data.ID,
		data.ParentVariantID,
		data.ChildVariantID,
		data.QuantityRatio,
		data.RepackCostPerUnit,
		data.RepackTimeMinutes,
		data.CreatedBy,
		data.UpdatedBy,
	)
	if err != nil {
		return err
	}
	return nil
}
