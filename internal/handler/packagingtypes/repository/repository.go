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

type PackagingType interface {
	InsertTransaction(ctx context.Context, data *PackagingTypeData) error
}
