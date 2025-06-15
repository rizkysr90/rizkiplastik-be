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
	updateSizeUnitSQL = `
		UPDATE size_units
		SET
			name = $2,
			unit_type = $3,
			description = $4,
			updated_by = $5,
			is_active = $6,
			updated_at = NOW()
		WHERE id = $1
	`
	findSizeUnitsPaginatedSQL = `
		SELECT
			id,
			name,
			code,
			unit_type,
			is_active,
			created_at,
			updated_at,
			COUNT(*) OVER () AS total_count
		FROM size_units
		WHERE 
			($1 = '' OR name LIKE '%' || $1 || '%') AND
			($2 = '' OR code = $2) AND
			($3 = '' OR unit_type = $3) AND
			(
				CASE
					WHEN $4 = 'TRUE' THEN is_active = true
					WHEN $4 = 'FALSE' THEN is_active = false
					ELSE TRUE
				END
			)
		ORDER BY created_at DESC
		LIMIT $5 OFFSET $6
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
func (s *SizeUnits) UpdateTrasaction(
	ctx context.Context,
	data repository.SizeUnitData) error {
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	rows, err := tx.Exec(ctx, updateSizeUnitSQL,
		data.SizeUnitID,
		data.SizeUnitName,
		data.SizeUnitType,
		data.SizeUnitDescription,
		data.UpdatedBy,
		data.IsActive,
	)
	if err != nil {
		return err
	}
	if rows.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (s *SizeUnits) FindPaginatedSizeUnits(
	ctx context.Context,
	filter repository.SizeUnitFilter) (
	[]repository.SizeUnitData, int, error) {
	rows, err := s.db.Query(ctx, findSizeUnitsPaginatedSQL,
		filter.SizeUnitName,
		filter.SizeUnitCode,
		filter.SizeUnitType,
		filter.IsActive,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var sizeUnits []repository.SizeUnitData
	var totalCount int
	for rows.Next() {
		var sizeUnit repository.SizeUnitData
		err = rows.Scan(
			&sizeUnit.SizeUnitID,
			&sizeUnit.SizeUnitName,
			&sizeUnit.SizeUnitCode,
			&sizeUnit.SizeUnitType,
			&sizeUnit.IsActive,
			&sizeUnit.CreatedAt,
			&sizeUnit.UpdatedAt,
			&totalCount,
		)
		if err != nil {
			return nil, 0, err
		}
		sizeUnits = append(sizeUnits, sizeUnit)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return sizeUnits, totalCount, nil
}
