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

type Category interface {
	InsertTransaction(ctx context.Context, data *CategoryData) error
}
