package products

import (
	"github.com/shopspring/decimal"
)

/*Base Data*/
type VariantObject struct {
	VariantName     *string             `json:"variant_name"`
	PackagingTypeID string              `json:"packaging_type_id"`
	SizeValue       float32             `json:"size_value"`
	SizeUnitID      string              `json:"size_unit_id"`
	CostPrice       *decimal.Decimal    `json:"cost_price"`
	SellPrice       decimal.Decimal     `json:"sell_price"`
	RepackRecipe    *RepackRecipeObject `json:"repack_recipe"`
}
type RepackRecipeObject struct {
	ParentVariantID   string          `json:"parent_variant_id"`
	QuantityRatio     float32         `json:"quantity_ratio"`
	RepackCostPerUnit decimal.Decimal `json:"repack_cost_per_unit"`
	RepackTimeMinutes int             `json:"repack_time_minutes"`
}
type Product struct {
	BaseName    string `json:"base_name"`
	ProductType string `json:"product_type"`
	CategoryID  string `json:"category_id"`
}
