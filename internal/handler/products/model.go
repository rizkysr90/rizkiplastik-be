package products

// CreateProductRequest represents data needed to create a product
type CreateProductRequest struct {
	Name                        string   `json:"name" binding:"required,max=50"`
	CostPrice                   float32  `json:"cost_price" binding:"required"`
	GrossProfitPercentage       float32  `json:"gross_profit_percentage" binding:"required"`
	VarianGrossProfitPercentage *float32 `json:"varian_gross_profit_percentage"`
	ShopeeCategory              string   `json:"shopee_category" binding:"required,oneof=A B C D E"`
	ShopeeVarianName            string   `json:"shopee_varian_name" binding:"max=100"`
	ShopeeName                  string   `json:"shopee_name" binding:"required,max=100"`
}

// UpdateProductRequest represents data needed to update a product
type UpdateProductRequest struct {
	Name                        string   `json:"name" binding:"required,max=50"`
	GrossProfitPercentage       float32  `json:"gross_profit_percentage" binding:"required"`
	VarianGrossProfitPercentage *float32 `json:"varian_gross_profit_percentage"`
	ShopeeCategory              string   `json:"shopee_category" binding:"required,oneof=A B C D E"`
	ShopeeVarianName            string   `json:"shopee_varian_name" binding:"max=100"`
	ShopeeName                  string   `json:"shopee_name" binding:"required,max=100"`
	CostPrice                   float32  `json:"cost_price" binding:"required"`
}

// PaginationResponse represents standard pagination metadata
type PaginationResponse struct {
	PageSize   int `json:"page_size"`
	PageNumber int `json:"page_number"`
}

// GetProductsResponse represents the response for the get products endpoint
type GetProductsResponse struct {
	Metadata PaginationResponse `json:"metadata"`
	Data     []Product          `json:"data"`
}

// GetProductResponse represents the response for the get product by ID endpoint
type GetProductResponse struct {
	Data Product `json:"data"`
}
