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
	findPaginatedPackagingTypesQuery = `
		SELECT 
			id, 
			name, 
			code, 
			is_active, 
			created_at, 
			updated_at, 
			COUNT(*) OVER () AS total_count
		FROM packaging_types
		WHERE $1 = '' OR name LIKE '%' || $1 || '%'
		AND ($2 = '' OR code = $2)
		AND (
			CASE 
				WHEN $3 = 'TRUE' THEN is_active = true 
				WHEN $3 = 'FALSE' THEN is_active = false 
				ELSE TRUE
			END
		)
		ORDER BY created_at DESC
		LIMIT $4 OFFSET $5
	`
	findByCategoryIDExtendedQuery = `
		SELECT 
			id, 
			name, 
			code,
			description,
			is_active,
			created_at,
			updated_at,
			created_by,
			updated_by
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

func (p *PackagingType) FindPaginatedPackagingTypes(
	ctx context.Context,
	filter *repository.PackagingTypeFilter) (
	[]repository.PackagingTypeData, int, error) {
	rows, err := p.db.Query(ctx, findPaginatedPackagingTypesQuery,
		filter.Name,
		filter.Code,
		filter.IsActive,
		filter.Limit,
		filter.Offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	var packagingTypes []repository.PackagingTypeData
	var totalCount int
	for rows.Next() {
		var packagingType repository.PackagingTypeData
		if err := rows.Scan(
			&packagingType.ID,
			&packagingType.Name,
			&packagingType.Code,
			&packagingType.IsActive,
			&packagingType.CreatedAt,
			&packagingType.UpdatedAt,
			&totalCount,
		); err != nil {
			return nil, 0, err
		}
		packagingTypes = append(packagingTypes, packagingType)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return packagingTypes, totalCount, nil
}
func (p *PackagingType) FindByCategoryIDExtended(
	ctx context.Context,
	categoryID string) (*repository.PackagingTypeData, error) {
	row := p.db.QueryRow(ctx, findByCategoryIDExtendedQuery, categoryID)
	var packagingType repository.PackagingTypeData
	if err := row.Scan(
		&packagingType.ID,
		&packagingType.Name,
		&packagingType.Code,
		&packagingType.Description,
		&packagingType.IsActive,
		&packagingType.CreatedAt,
		&packagingType.UpdatedAt,
		&packagingType.CreatedBy,
		&packagingType.UpdatedBy,
	); err != nil {
		if err == pgx.ErrNoRows {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return &packagingType, nil
}
