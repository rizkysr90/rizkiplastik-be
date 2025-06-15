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

type VarianType interface {
	InsertTransaction(
		ctx context.Context, data *VarianTypeData) error
	UpdateTransaction(
		ctx context.Context, data *VarianTypeData) error
}
