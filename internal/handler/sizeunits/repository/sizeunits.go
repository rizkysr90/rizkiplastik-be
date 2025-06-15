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

type SizeUnitFilter struct {
	SizeUnitName string
	SizeUnitCode string
	SizeUnitType string
	IsActive     string
	Limit        int
	Offset       int
}

type SizeUnits interface {
	InsertTransaction(ctx context.Context, data SizeUnitData) error
	UpdateTrasaction(ctx context.Context, data SizeUnitData) error
	FindPaginatedSizeUnits(
		ctx context.Context,
		filter SizeUnitFilter) (
		[]SizeUnitData, int, error)
}
