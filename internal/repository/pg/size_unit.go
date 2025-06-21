package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrSizeUnitNotFound = errors.New("size unit not found")
)

type SizeUnit struct {
	db *pgxpool.Pool
}

func NewSizeUnit(db *pgxpool.Pool) *SizeUnit {
	return &SizeUnit{db: db}
}
func (pg *SizeUnit) checkSizeUnitID(
	ctx context.Context,
	tx pgx.Tx,
	inputSizeUnitID string,
) error {
	var sizeUnitID string
	err := tx.QueryRow(
		ctx,
		checkSizeUnitIDSQL,
		inputSizeUnitID,
	).Scan(&sizeUnitID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrSizeUnitNotFound
		}
		return err
	}
	return nil
}
