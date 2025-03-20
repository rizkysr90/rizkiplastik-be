package products

import (
	"time"

	"github.com/google/uuid"
)

// Product represents a product entity
type Product struct {
	ID                    uuid.UUID  `json:"id"`
	Name                  string     `json:"name"`
	CostPrice             float32    `json:"cost_price"`
	GrossProfitPercentage float32    `json:"gross_profit_percentage"`
	ShopeeCategory        string     `json:"shopee_category"`
	ShopeeSalePrice       float32    `json:"shopee_sale_price,omitempty"`
	ShopeeFee             float32    `json:"shopee_fee,omitempty"`
	CreatedAt             time.Time  `json:"-"`
	UpdatedAt             *time.Time `json:"-"`
	DeletedAt             *time.Time `json:"-"`
}
