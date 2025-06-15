package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/repository"
)

type SizeUnits struct {
	db *pgxpool.Pool
}

func NewSizeUnits(db *pgxpool.Pool) *SizeUnits {
	return &SizeUnits{db: db}
}

const (
	findSizeUnitByCodeSQL = `
		SELECT
			id,
			code
		FROM size_units
		WHERE code = $1
	`
	insertSizeUnitSQL = `
		INSERT INTO size_units (
			id,
			name,
			code,
			unit_type,
			description,
			created_by,
			updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7)
	`
)

func (s *SizeUnits) InsertTransaction(
	ctx context.Context, data repository.SizeUnitData) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var sizeUnitByCode repository.SizeUnitData
	err = tx.QueryRow(ctx, findSizeUnitByCodeSQL, data.SizeUnitCode).Scan(
		&sizeUnitByCode.SizeUnitID,
		&sizeUnitByCode.SizeUnitCode,
	)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	if sizeUnitByCode.SizeUnitID != "" {
		return constants.ErrAlreadyExists
	}
	_, err = tx.Exec(ctx, insertSizeUnitSQL,
		data.SizeUnitID,
		data.SizeUnitName,
		data.SizeUnitCode,
		data.SizeUnitType,
		data.SizeUnitDescription,
		data.CreatedBy,
		data.UpdatedBy,
	)
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
