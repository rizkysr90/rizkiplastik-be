package onlinetransactions

import (
	"time"

	"github.com/google/uuid"
)

// OnlineTransaction represents an online transaction entity
type OnlineTransaction struct {
	ID              uuid.UUID                  `json:"id"`
	Type            string                     `json:"type"` // SHOPEE, LAZADA, TOKOPEDIA, TIKTOK
	OrderNumber     string                     `json:"order_number"`
	CreatedDate     time.Time                  `json:"created_date"`
	PeriodMonth     int                        `json:"period_month"`
	PeriodYear      int                        `json:"period_year"`
	TotalBaseAmount float32                    `json:"total_base_amount"`
	TotalSaleAmount float32                    `json:"total_sale_amount"`
	TotalNetProfit  float32                    `json:"total_net_profit"`
	TotalFeeAmount  float32                    `json:"total_fee_amount"`
	CreatedBy       string                     `json:"created_by"`
	Products        []OnlineTransactionProduct `json:"products,omitempty"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       *time.Time                 `json:"-"`
	DeletedAt       *time.Time                 `json:"-"`
}

// OnlineTransactionProduct represents a product in an online transaction
type OnlineTransactionProduct struct {
	ID          uuid.UUID `json:"id"`
	ProductName string    `json:"product_name"`
	CostPrice   float32   `json:"cost_price"`
	SalePrice   float32   `json:"sale_price"`
	Quantity    int       `json:"quantity"`
	FeeAmount   float32   `json:"fee_amount"`
}
