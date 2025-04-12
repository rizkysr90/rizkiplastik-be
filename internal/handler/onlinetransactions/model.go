package onlinetransactions

// CreateOnlineTransactionRequest represents data needed to create an online transaction
type CreateOnlineTransactionRequest struct {
	Type        string                     `json:"type" binding:"required,oneof=SHOPEE LAZADA TOKOPEDIA TIKTOK"`
	OrderNumber string                     `json:"order_number" binding:"required"`
	CreatedDate string                     `json:"created_date" binding:"required"` // Format: YYYY-MM-DD
	Products    []CreateTransactionProduct `json:"products" binding:"required,dive,required"`
}

// CreateTransactionProduct represents a product in the create transaction request
type CreateTransactionProduct struct {
	ProductID string `json:"product_id" binding:"required,max=50"`
	Quantity  int    `json:"quantity" binding:"required,min=1"`
}

// UpdateOnlineTransactionRequest represents data needed to update an online transaction
type UpdateOnlineTransactionRequest struct {
	Type        string                     `json:"type" binding:"required,oneof=SHOPEE LAZADA TOKOPEDIA TIKTOK"`
	OrderNumber string                     `json:"order_number" binding:"required"`
	CreatedDate string                     `json:"created_date" binding:"required"` // Format: YYYY-MM-DD
	Products    []CreateTransactionProduct `json:"products" binding:"required,dive,required"`
}

// PaginationResponse represents standard pagination metadata
type PaginationResponse struct {
	PageSize   int `json:"page_size"`
	PageNumber int `json:"page_number"`
	TotalCount int `json:"total_count"`
	TotalPages int `json:"total_pages"`
}

// GetOnlineTransactionsResponse represents the response for the get transactions endpoint
type GetOnlineTransactionsResponse struct {
	Metadata PaginationResponse  `json:"metadata"`
	Data     []OnlineTransaction `json:"data"`
}

// GetOnlineTransactionResponse represents the response for the get transaction by ID endpoint
type GetOnlineTransactionResponse struct {
	Data OnlineTransaction `json:"data"`
}
