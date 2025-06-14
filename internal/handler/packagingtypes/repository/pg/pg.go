package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/repository"
)

type PackagingType struct {
	db *pgxpool.Pool
}

func NewPackagingType(db *pgxpool.Pool) *PackagingType {
	return &PackagingType{db: db}
}

const (
	insertPackagingTypeQuery = `
		INSERT INTO packaging_types (
			id, name, code, description, created_by, updated_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`
	findPackagingTypeByCodeQuery = `
		SELECT id, name, code
		FROM packaging_types
		WHERE code = $1
	`
	updatePackagingTypeQuery = `
		UPDATE packaging_types
		SET name = $2, code = $3, description = $4, is_active = $5, 
		updated_by = $6, updated_at = NOW()
		WHERE id = $1
	`
	findPackagingTypeByIDQuery = `
		SELECT id, name, code
		FROM packaging_types
		WHERE id = $1
	`
)

func (p *PackagingType) getByCode(ctx context.Context,
	tx pgx.Tx, code string) (*repository.PackagingTypeData, error) {
	row := tx.QueryRow(ctx, findPackagingTypeByCodeQuery, code)
	var packagingType repository.PackagingTypeData
	if err := row.Scan(
		&packagingType.ID,
		&packagingType.Name,
		&packagingType.Code,
	); err != nil && err != pgx.ErrNoRows {
		return nil, err
	}
	return &packagingType, nil
}
func (p *PackagingType) InsertTransaction(
	ctx context.Context, data *repository.PackagingTypeData) error {
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	packagingType, err := p.getByCode(ctx, tx, data.Code)
	if err != nil {
		return err
	}
	if packagingType.ID != "" {
		return constants.ErrAlreadyExists
	}
	_, err = tx.Exec(ctx,
		insertPackagingTypeQuery,
		data.ID,
		data.Name,
		data.Code,
		data.Description,
		data.CreatedBy,
		data.UpdatedBy,
	)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}

func (p *PackagingType) UpdateTransaction(
	ctx context.Context, data *repository.PackagingTypeData) error {
	tx, err := p.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	row := tx.QueryRow(ctx, findPackagingTypeByIDQuery, data.ID)
	var packagingType repository.PackagingTypeData
	if err := row.Scan(
		&packagingType.ID,
		&packagingType.Name,
		&packagingType.Code,
	); err != nil {
		if err == pgx.ErrNoRows {
			return pgx.ErrNoRows
		}
		return err
	}
	if packagingType.Code != data.Code {
		packagingTypeByCode, err := p.getByCode(ctx, tx, data.Code)
		if err != nil {
			return err
		}
		if packagingTypeByCode.ID != "" {
			return constants.ErrAlreadyExists
		}
	}
	_, err = tx.Exec(
		ctx, updatePackagingTypeQuery,
		data.ID,
		data.Name,
		data.Code,
		data.Description,
		data.IsActive,
		data.UpdatedBy,
	)
	if err != nil {
		return err
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
