package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type ProductType string

const (
	ProductTypeRepack  ProductType = "REPACK"
	ProductTypeVariant ProductType = "VARIANT"
	ProductTypeSingle  ProductType = "SINGLE"
)

type ProductData struct {
	ID          string
	BaseName    string
	CategoryID  string
	ProductType ProductType
	CreatedBy   string
	UpdatedBy   string
}

type ProductRepository interface {
	InsertTransaction(
		ctx context.Context,
		tx pgx.Tx,
		data *ProductData,
	) error
}
