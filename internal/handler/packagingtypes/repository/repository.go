package repository

import (
	"context"
	"database/sql"
	"time"
)

type PackagingTypeData struct {
	ID          string
	Name        string
	Code        string
	Description sql.NullString
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CreatedBy   string
	UpdatedBy   string
}
type PackagingTypeFilter struct {
	Name     string
	Code     string
	IsActive string
	Limit    int
	Offset   int
}
type PackagingType interface {
	InsertTransaction(ctx context.Context, data *PackagingTypeData) error
	UpdateTransaction(
		ctx context.Context, data *PackagingTypeData) error
	FindPaginatedPackagingTypes(
		ctx context.Context,
		filter *PackagingTypeFilter) (
		[]PackagingTypeData, int, error)
}
