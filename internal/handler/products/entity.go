package products

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product entity
type Product struct {
	ID                          uuid.UUID  `json:"id"`
	Name                        string     `json:"name"`
	CostPrice                   float32    `json:"cost_price"`
	GrossProfitPercentage       float32    `json:"gross_profit_percentage"`
	VarianGrossProfitPercentage float32    `json:"varian_gross_profit_percentage"`
	ShopeeFreeDeliveryFee       float32    `json:"shopee_fee_free_delivery_fee"`
	ShopeeCategory              string     `json:"shopee_category"`
	ShopeeVarianName            string     `json:"shopee_varian_name"`
	ShopeeName                  string     `json:"shopee_name"`
	ShopeeSalePrice             float32    `json:"shopee_sale_price,omitempty"`
	ShopeeFee                   float32    `json:"shopee_fee,omitempty"`
	CreatedAt                   time.Time  `json:"-"`
	UpdatedAt                   *time.Time `json:"-"`
	DeletedAt                   *time.Time `json:"-"`
}
