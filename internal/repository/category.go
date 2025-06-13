package repository

import (
	"context"
	"time"
)

// CategoryData is the data structure for the category table
type CategoryData struct {
	ID          string
	Name        string
	Code        string
	Description string
	IsActive    bool
	CreatedBy   string
	UpdatedBy   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// CategoryData Filter
type CategoryDataFilter struct {
	CategoryName string
	CategoryCode string
	IsActive     string
	PageSize     int
	Offset       int
}
type Category interface {
	InsertTransaction(ctx context.Context, data *CategoryData) error
	Update(ctx context.Context, data *CategoryData) error
	GetList(ctx context.Context, filter *CategoryDataFilter) ([]CategoryData, int, error)
}
