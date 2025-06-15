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
	findVarianTypeByIdQuery = `
		SELECT id, name
		FROM variant_types WHERE id = $1
	`
	updateVariantTypeQuery = `
		UPDATE variant_types SET
			name = $2,
			description = $3,
			is_active = $4,
			updated_by = $5,
			updated_at = NOW()
		WHERE id = $1
	`
)

func (v *VarianTypes) findVarianTypeByName(
	ctx context.Context, tx pgx.Tx, name string) (*repository.VarianTypeData, error) {
	var varianTypeByName repository.VarianTypeData
	err := tx.QueryRow(ctx, findVarianTypeByNameQuery, name).Scan(
		&varianTypeByName.ID,
		&varianTypeByName.Name,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if varianTypeByName.ID != "" {
		return nil, constants.ErrAlreadyExists
	}
	return &varianTypeByName, nil
}
func (v *VarianTypes) InsertTransaction(
	ctx context.Context, data *repository.VarianTypeData) error {
	tx, err := v.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = v.findVarianTypeByName(ctx, tx, data.Name)
	if err != nil {
		return err
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
func (v *VarianTypes) UpdateTransaction(
	ctx context.Context, data *repository.VarianTypeData) error {
	tx, err := v.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var varianTypeById repository.VarianTypeData
	err = tx.QueryRow(ctx, findVarianTypeByIdQuery, data.ID).Scan(
		&varianTypeById.ID,
		&varianTypeById.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return err
	}
	if varianTypeById.Name != data.Name {
		_, err = v.findVarianTypeByName(ctx, tx, data.Name)
		if err != nil {
			return err
		}
	}
	_, err = tx.Exec(ctx, updateVariantTypeQuery,
		data.ID,
		data.Name,
		data.Description,
		data.IsActive,
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
