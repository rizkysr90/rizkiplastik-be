package onlinetransactions

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/products"
	"github.com/rizkysr90/rizkiplastik-be/internal/middleware"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

// OnlineTransactions handles HTTP requests for products
type Handler struct {
	db *pgxpool.Pool
}

// OnlineTransactions creates a new product handler
func NewOnlineTransactions(db *pgxpool.Pool) *Handler {
	return &Handler{db: db}
}

// RegisterRoutes registers all online transaction related routes
func (h *Handler) RegisterRoutes(
	router *gin.Engine,
	authMiddleware *middleware.AuthMiddleware) {

	// Admin routes (only admin can access)
	admin := router.Group("/api/v1/online-transactions")
	admin.Use(authMiddleware.RequireAuth(), authMiddleware.RequireRole("ADMIN"))
	{
		admin.GET("", h.GetOnlineTransactions)
		admin.GET("/:id", h.GetOnlineTransaction)
		admin.POST("", h.CreateOnlineTransaction)
		// admin.PUT("/:id", h.UpdateOnlineTransaction)
		admin.DELETE("/:id", h.DeleteOnlineTransaction)
	}
}

// CreateOnlineTransaction handles the creation of a new online transaction
func (h *Handler) CreateOnlineTransaction(c *gin.Context) {
	var req CreateOnlineTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get username from context (set by auth middleware)
	username, ok := util.GetUsernameFromContext(c)
	if !ok {
		// Handle the error case, e.g., user not authenticated
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	// Begin transaction
	tx, err := h.db.Begin(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to begin transaction"})
		return
	}
	defer tx.Rollback(c)
	listProductID := []string{}
	listProduct := []products.Product{}
	mapProductIDWithQty := map[string]int{}
	mapProductIDWithSalePrice := map[string]float32{}
	mapProductIDWithFee := map[string]float32{}

	for _, product := range req.Products {
		listProductID = append(listProductID, product.ProductID)
		mapProductIDWithQty[product.ProductID] = product.Quantity
	}
	rows, err := tx.Query(c, ` 
		SELECT id, name, cost_price, gross_profit_percentage, shopee_category
        FROM products 
        WHERE id = ANY($1::uuid[])
		`, listProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid get detail product"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var temp products.Product
		err = rows.Scan(
			&temp.ID,
			&temp.Name,
			&temp.CostPrice,
			&temp.GrossProfitPercentage,
			&temp.ShopeeCategory,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid scan"})
			return
		}
		listProduct = append(listProduct, temp)
	}
	err = rows.Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid rows err"})
		return
	}
	// Create the transaction
	txID := uuid.New()
	now := time.Now().UTC()

	// Parse created_date from request
	createdDate, err := time.Parse("2006-01-02", req.CreatedDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
		return
	}

	// Calculate totals
	var totalBaseAmount, totalSaleAmount, totalNetProfit, totalFeeAmount float32 = 0, 0, 0, 0
	for _, product := range listProduct {
		productSalePrice, productFee := products.CalculateShopeePricing(product.CostPrice, product.GrossProfitPercentage, product.ShopeeCategory)
		mapProductIDWithSalePrice[product.ID.String()] = productSalePrice
		mapProductIDWithFee[product.ID.String()] = productFee
		totalBaseAmount += product.CostPrice * float32(mapProductIDWithQty[product.ID.String()])
		totalSaleAmount += productSalePrice * float32(mapProductIDWithQty[product.ID.String()])
		totalFeeAmount += productFee * float32(mapProductIDWithQty[product.ID.String()])
	}
	totalNetProfit = (totalSaleAmount - totalFeeAmount) - totalBaseAmount

	// Extract month and year from created date
	periodMonth := int(createdDate.Month())
	periodYear := createdDate.Year()

	// Insert transaction
	transactionQuery := `
		INSERT INTO online_transactions (
			id, type, order_number, created_date, period_month, period_year,
			total_base_amount, total_sale_amount, total_net_profit, created_by,
			created_at, total_fee_amount
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`

	_, err = tx.Exec(c, transactionQuery,
		txID,
		strings.ToUpper(req.Type),
		strings.ToUpper(req.OrderNumber),
		createdDate,
		periodMonth,
		periodYear,
		totalBaseAmount,
		totalSaleAmount,
		totalNetProfit,
		username,
		now,
		totalFeeAmount,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction: " + err.Error()})
		return
	}

	// Insert products
	for _, product := range listProduct {
		productQuery := `
			INSERT INTO online_transaction_products (
				id, online_transaction_id, product_name, cost_price, sale_price, quantity, fee_amount
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
		`

		_, err = tx.Exec(c, productQuery,
			uuid.New(),
			txID,
			product.Name,
			product.CostPrice,
			mapProductIDWithSalePrice[product.ID.String()],
			mapProductIDWithQty[product.ID.String()],
			mapProductIDWithFee[product.ID.String()],
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create transaction products: " + err.Error()})
			return
		}
	}

	// Commit transaction
	if err := tx.Commit(c); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{})
}

// DeleteOnlineTransaction handles soft deletion of a transaction
func (h *Handler) DeleteOnlineTransaction(c *gin.Context) {
	id := c.Param("id")
	transactionID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	now := time.Now().UTC()
	query := `
		UPDATE online_transactions
		SET deleted_at = $1
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := h.db.Exec(c, query, now, transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete transaction"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// GetOnlineTransaction retrieves a transaction by ID with its products
func (h *Handler) GetOnlineTransaction(c *gin.Context) {
	id := c.Param("id")
	transactionID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	// Get transaction data
	var transaction OnlineTransaction
	transactionQuery := `
		SELECT id, type, order_number, created_date, period_month, period_year,
			   total_base_amount, total_sale_amount, total_net_profit, created_by,
			   created_at, total_fee_amount
		FROM online_transactions
		WHERE id = $1 AND deleted_at IS NULL
	`

	err = h.db.QueryRow(c, transactionQuery, transactionID).Scan(
		&transaction.ID,
		&transaction.Type,
		&transaction.OrderNumber,
		&transaction.CreatedDate,
		&transaction.PeriodMonth,
		&transaction.PeriodYear,
		&transaction.TotalBaseAmount,
		&transaction.TotalSaleAmount,
		&transaction.TotalNetProfit,
		&transaction.CreatedBy,
		&transaction.CreatedAt,
		&transaction.TotalFeeAmount,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transaction"})
		return
	}

	// Get transaction products
	productsQuery := `
		SELECT id, product_name, cost_price, sale_price, quantity, fee_amount
		FROM online_transaction_products
		WHERE online_transaction_id = $1
	`

	rows, err := h.db.Query(c, productsQuery, transactionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transaction products"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product OnlineTransactionProduct
		err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.CostPrice,
			&product.SalePrice,
			&product.Quantity,
			&product.FeeAmount,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product data"})
			return
		}
		transaction.Products = append(transaction.Products, product)
	}

	response := GetOnlineTransactionResponse{
		Data: transaction,
	}

	c.JSON(http.StatusOK, response)
}

// GetOnlineTransactions retrieves a list of transactions with pagination and filtering
func (h *Handler) GetOnlineTransactions(c *gin.Context) {

	pageSize := 50
	pageNumber := 0
	typeFilter := ""
	orderNumberFilter := ""
	startDate := ""
	endDate := ""

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

	// Parse filters
	typeFilter = c.Query("type")
	orderNumberFilter = c.Query("order_number")
	startDate = c.Query("start_date")
	endDate = c.Query("end_date")

	// Calculate offset
	offset := pageNumber * pageSize
	// Build query with conditionals for filters
	query := `
		SELECT id, type, order_number, created_date, period_month, period_year,
			   total_base_amount, total_sale_amount, total_net_profit, created_by,
			   created_at, COUNT(*) OVER() AS total_count, total_fee_amount
		FROM online_transactions
		WHERE 

		deleted_at IS NULL
	`
	// Add filters conditionally
	params := []interface{}{}
	paramCounter := 1
	if orderNumberFilter != "" {
		query += ` AND order_number = $` + strconv.Itoa(paramCounter)
		params = append(params, orderNumberFilter)
		paramCounter++
	}
	if typeFilter != "" {
		query += ` AND type = $` + strconv.Itoa(paramCounter)
		params = append(params, strings.ToUpper(typeFilter))
		paramCounter++
	}
	if startDate != "" {
		query += ` AND created_date >= $` + strconv.Itoa(paramCounter)
		params = append(params, startDate)
		paramCounter++
	}
	if endDate != "" {
		query += ` AND created_date <= $` + strconv.Itoa(paramCounter)
		params = append(params, endDate)
		paramCounter++
	}
	// Add ordering, limit and offset
	query += `
		ORDER BY created_date DESC
		LIMIT $` + strconv.Itoa(paramCounter) + ` OFFSET $` + strconv.Itoa(paramCounter+1)
	params = append(params, pageSize, offset)

	rows, err := h.db.Query(c, query, params...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve transactions: " + err.Error()})
		return
	}
	defer rows.Close()

	transactions := []OnlineTransaction{}
	var totalCount int64 = 0

	for rows.Next() {
		var transaction OnlineTransaction
		err := rows.Scan(
			&transaction.ID,
			&transaction.Type,
			&transaction.OrderNumber,
			&transaction.CreatedDate,
			&transaction.PeriodMonth,
			&transaction.PeriodYear,
			&transaction.TotalBaseAmount,
			&transaction.TotalSaleAmount,
			&transaction.TotalNetProfit,
			&transaction.CreatedBy,
			&transaction.CreatedAt,
			&totalCount,
			&transaction.TotalFeeAmount,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan transaction data"})
			return
		}

		transactions = append(transactions, transaction)
	}

	response := GetOnlineTransactionsResponse{
		Metadata: PaginationResponse{
			PageSize:   pageSize,
			PageNumber: pageNumber,
			TotalCount: int(totalCount),
			TotalPages: int(totalCount+int64(pageSize)-1) / pageSize,
		},
		Data: transactions,
	}

	c.JSON(http.StatusOK, response)

}
