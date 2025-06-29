package products

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/middleware"
)

// ProductHandler handles HTTP requests for products
type ProductHandler struct {
	db *pgxpool.Pool
}

// NewProductHandler creates a new product handler
func NewProductHandler(db *pgxpool.Pool) *ProductHandler {
	return &ProductHandler{db: db}
}

// RegisterRoutes registers all product-related routes
func (h *ProductHandler) RegisterRoutes(
	router *gin.Engine,
	authMiddleware *middleware.AuthMiddleware) {
	private := router.Group("api/v1/products")
	private.Use(authMiddleware.RequireAuth())
	{
		private.GET("", h.GetProducts)    // Anyone can list products
		private.GET("/:id", h.GetProduct) // Anyone can view a product
	}
	// Admin routes (only admin can access)
	admin := router.Group("/api/v1/products")
	admin.Use(authMiddleware.RequireAuth(), authMiddleware.RequireRole("ADMIN"))
	{
		admin.POST("", h.CreateProduct)
		admin.PUT("/:id", h.UpdateProduct)
		admin.DELETE("/:id", h.DeleteProduct)
	}
}

// CreateProduct handles the creation of a new product
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var varianGPP interface{}
	if req.VarianGrossProfitPercentage != nil && *req.VarianGrossProfitPercentage != 0 {
		varianGPP = *req.VarianGrossProfitPercentage
	} else {
		varianGPP = nil
	}

	// Handle empty variant name
	var shopeeVarianName interface{}
	if req.ShopeeVarianName == "" {
		shopeeVarianName = nil
	} else {
		shopeeVarianName = strings.ToUpper(req.ShopeeVarianName)
	}

	product := Product{
		ID:                    uuid.New(),
		Name:                  req.Name,
		CostPrice:             req.CostPrice,
		GrossProfitPercentage: req.GrossProfitPercentage,
		ShopeeCategory:        req.ShopeeCategory,
		ShopeeName:            req.ShopeeName,
		CreatedAt:             time.Now().UTC(),
	}

	query := `
		INSERT INTO products (
		id, name, cost_price, gross_profit_percentage, varian_gross_profit_percentage, shopee_category, shopee_varian_name, shopee_name, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := h.db.Exec(c, query,
		product.ID,
		strings.ToUpper(product.Name),
		product.CostPrice,
		product.GrossProfitPercentage,
		varianGPP,
		strings.ToUpper(product.ShopeeCategory),
		shopeeVarianName,
		strings.ToUpper(product.ShopeeName),
		product.CreatedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// UpdateProduct handles updating an existing product
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now().UTC()
	var varianGPP interface{}
	if req.VarianGrossProfitPercentage != nil && *req.VarianGrossProfitPercentage != 0 {
		varianGPP = *req.VarianGrossProfitPercentage
	} else {
		varianGPP = nil
	}

	// Handle empty variant name
	var shopeeVarianName interface{}
	if req.ShopeeVarianName == "" {
		shopeeVarianName = nil
	} else {
		shopeeVarianName = strings.TrimSpace(strings.ToUpper(req.ShopeeVarianName))
	}

	query := `
		UPDATE products
		SET name = $1, gross_profit_percentage = $2, varian_gross_profit_percentage = $3, shopee_category = $4, shopee_varian_name = $5, shopee_name = $6, updated_at = $7,
		cost_price = $8
		WHERE id = $9 AND deleted_at IS NULL
	`

	result, err := h.db.Exec(c, query,
		strings.TrimSpace(strings.ToUpper(req.Name)),
		req.GrossProfitPercentage,
		varianGPP,
		strings.ToUpper(req.ShopeeCategory),
		shopeeVarianName,
		strings.TrimSpace(strings.ToUpper(req.ShopeeName)),
		now,
		req.CostPrice,
		productID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// DeleteProduct handles soft deletion of a product
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	now := time.Now().UTC()
	query := `
		UPDATE products
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := h.db.Exec(c, query, now, productID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// GetProduct retrieves a product by ID
func (h *ProductHandler) GetProduct(c *gin.Context) {
	id := c.Param("id")
	productID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var product Product
	query := `
		SELECT id,
		name, 
		cost_price, 
		gross_profit_percentage, 
		COALESCE(varian_gross_profit_percentage, 0), 
		shopee_category, 
		COALESCE(shopee_varian_name, ''), 
		COALESCE(shopee_name, ''),
		shopee_fee_free_delivery_fee
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`

	var varianGPP sql.NullFloat64
	err = h.db.QueryRow(c, query, productID).Scan(
		&product.ID,
		&product.Name,
		&product.CostPrice,
		&product.GrossProfitPercentage,
		&varianGPP,
		&product.ShopeeCategory,
		&product.ShopeeVarianName,
		&product.ShopeeName,
		&product.ShopeeFreeDeliveryFee,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve product"})
		return
	}
	if varianGPP.Valid {
		product.VarianGrossProfitPercentage = float32(varianGPP.Float64)
	} else {
		product.VarianGrossProfitPercentage = 0
	}

	// Calculate shopee sale price and fee
	product.ShopeeSalePrice, product.ShopeeFee = CalculateShopeePricing(
		product.CostPrice,
		product.GrossProfitPercentage,
		product.ShopeeFreeDeliveryFee,
		product.ShopeeCategory,
	)

	response := GetProductResponse{
		Data: product,
	}

	c.JSON(http.StatusOK, response)
}

// GetProducts retrieves a list of products with pagination and filtering
func (h *ProductHandler) GetProducts(c *gin.Context) {
	pageSize := 50
	pageNumber := 0
	nameFilter := ""

	// Parse pagination parameters
	if pageSizeParam := c.Query("page_size"); pageSizeParam != "" {
		if size, err := strconv.Atoi(pageSizeParam); err == nil {
			pageSize = size
		}
	}

	if pageNumberParam := c.Query("page_number"); pageNumberParam != "" {
		if page, err := strconv.Atoi(pageNumberParam); err == nil {
			pageNumber = page
		}
	}

	// Parse name filter
	nameFilter = c.Query("name")
	if nameFilter == "" {
		nameFilter = "" // Ensure it's an empty string if not provided
	}

	// Calculate offset
	offset := pageNumber * pageSize

	// Use a single query with window function for count and data
	query := `
		SELECT 
			id, 
			name, 
			cost_price, 
			gross_profit_percentage, 
			varian_gross_profit_percentage, 
			shopee_category, 
			COALESCE(shopee_varian_name, ''), 
			COALESCE(shopee_name, ''), 
			shopee_fee_free_delivery_fee,
			COUNT(*) OVER() AS total_count
		FROM products
		WHERE deleted_at IS NULL AND
		($1 = '' OR name ILIKE '%' || $1 || '%')
		ORDER BY name
		LIMIT $2 OFFSET $3
	`

	rows, err := h.db.Query(c, query, nameFilter, pageSize, offset)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve products"})
		return
	}
	defer rows.Close()

	products := []Product{}
	var totalCount int64 = 0

	for rows.Next() {
		var product Product
		var varianGPP sql.NullFloat64
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.CostPrice,
			&product.GrossProfitPercentage,
			&varianGPP,
			&product.ShopeeCategory,
			&product.ShopeeVarianName,
			&product.ShopeeName,
			&product.ShopeeFreeDeliveryFee,
			&totalCount,
		)
		if err != nil {
			errMessage := fmt.Sprintf("Failed to scan product data: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
			return
		}

		if varianGPP.Valid {
			product.VarianGrossProfitPercentage = float32(varianGPP.Float64)
		} else {
			product.VarianGrossProfitPercentage = 0
		}

		// Calculate shopee sale price and fee
		product.ShopeeSalePrice, product.ShopeeFee = CalculateShopeePricing(
			product.CostPrice,
			product.GrossProfitPercentage,
			product.ShopeeFreeDeliveryFee,
			product.ShopeeCategory,
		)

		products = append(products, product)
	}

	response := GetProductsResponse{
		Metadata: PaginationResponse{
			PageSize:   pageSize,
			PageNumber: pageNumber,
		},
		Data: products,
	}

	c.JSON(http.StatusOK, response)
}
