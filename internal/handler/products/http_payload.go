package products

import "github.com/shopspring/decimal"

type CreateProductRequest struct {
	Product  Product         `json:"product"`
	Variants []VariantObject `json:"variants"`
}

type UpdateSingleProductTypeRequest struct {
	ProductID       string           `json:"product_id"`
	BaseName        string           `json:"base_name"`
	CategoryID      string           `json:"category_id"`
	PackagingTypeID string           `json:"packaging_type_id"`
	SizeValue       float32          `json:"size_value"`
	SizeUnitID      string           `json:"size_unit_id"`
	CostPrice       *decimal.Decimal `json:"cost_price"`
	SellPrice       decimal.Decimal  `json:"sell_price"`
}
