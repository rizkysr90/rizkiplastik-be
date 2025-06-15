package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/repository"
)

type VarianTypes struct {
	db *pgxpool.Pool
}

func NewVarianTypes(db *pgxpool.Pool) *VarianTypes {
	return &VarianTypes{db: db}
}

const (
	insertVariantTypeQuery = `
		INSERT INTO variant_types (
			id,
			name, 
			description, 
			created_by, 
			updated_by, 
			created_at, 
			updated_at
		) VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`
	findVarianTypeByNameQuery = `
		SELECT id, name
		FROM variant_types WHERE name = $1
	`
)

func (v *VarianTypes) InsertTransaction(
	ctx context.Context, data *repository.VarianTypeData) error {
	tx, err := v.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var varianTypeByName repository.VarianTypeData
	err = tx.QueryRow(ctx, findVarianTypeByNameQuery, data.Name).Scan(
		&varianTypeByName.ID,
		&varianTypeByName.Name,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if varianTypeByName.ID != "" {
		return constants.ErrAlreadyExists
	}

	_, err = tx.Exec(ctx, insertVariantTypeQuery,
		data.ID,
		data.Name,
		data.Description,
		data.CreatedBy,
		data.UpdatedBy,
	)
	if err != nil {
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
