package repository

import (
	"context"
	"database/sql"
	"time"
)

type VarianTypeData struct {
	ID          string
	Name        string
	Description sql.NullString
	IsActive    bool
	CreatedBy   string
	UpdatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
type VarianTypeFilter struct {
	Name     string
	IsActive string
	Offset   int
	Limit    int
}
type VarianType interface {
	InsertTransaction(
		ctx context.Context, data *VarianTypeData) error
	UpdateTransaction(
		ctx context.Context, data *VarianTypeData) error
	FindVarianTypePaginated(
		ctx context.Context,
		filter VarianTypeFilter) ([]VarianTypeData, int, error)
	FindVarianTypeByIdExtended(
		ctx context.Context, id string) (*VarianTypeData, error)
}
