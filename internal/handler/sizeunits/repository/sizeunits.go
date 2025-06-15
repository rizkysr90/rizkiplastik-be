package repository

import (
	"context"
	"database/sql"
	"time"
)

type SizeUnitData struct {
	SizeUnitID          string
	SizeUnitName        string
	SizeUnitCode        string
	SizeUnitType        string
	CreatedBy           string
	UpdatedBy           string
	SizeUnitDescription sql.NullString
	IsActive            bool
	CreatedAt           time.Time
	UpdatedAt           time.Time
}

type SizeUnits interface {
	InsertTransaction(ctx context.Context, data SizeUnitData) error
	UpdateTrasaction(ctx context.Context, data SizeUnitData) error
}
