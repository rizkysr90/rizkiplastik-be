package repository

type ProductType string

const (
	ProductTypeRepack  ProductType = "REPACK"
	ProductTypeVariant ProductType = "VARIANT"
)

type ProductData struct {
}

type ProductRepository interface {
}
