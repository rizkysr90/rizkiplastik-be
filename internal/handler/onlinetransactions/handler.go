package onlinetransactions

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
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
	"github.com/xuri/excelize/v2"
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
		admin.POST("/auto-input-excel", h.AutoInputOnlineTransactions)
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

type TransactionTotals struct {
	TotalBaseAmount float32
	TotalSaleAmount float32
	TotalFeeAmount  float32
	TotalNetProfit  float32
	Products        []TransactionProduct
}

type TransactionProduct struct {
	ProductName string
	CostPrice   float32
	SalePrice   float32
	Quantity    int
	FeeAmount   float32
}

func calculateTransactionTotals(products []TransactionProduct) TransactionTotals {
	var totals TransactionTotals
	for _, product := range products {
		totals.TotalBaseAmount += product.CostPrice * float32(product.Quantity)
		totals.TotalSaleAmount += product.SalePrice * float32(product.Quantity)
		totals.TotalFeeAmount += product.FeeAmount * float32(product.Quantity)
	}
	totals.TotalNetProfit = (totals.TotalSaleAmount - totals.TotalFeeAmount) - totals.TotalBaseAmount
	totals.Products = products
	return totals
}

// AutoInputOnlineTransactions handles Excel upload and auto input
func (h *Handler) AutoInputOnlineTransactions(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to open file"})
		return
	}
	defer f.Close()

	xl, err := excelize.OpenReader(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read Excel file"})
		return
	}
	sheetName := xl.GetSheetName(0)
	rows, err := xl.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read rows from Excel"})
		return
	}

	// Find column indexes
	colIdx := map[string]int{}
	for i, col := range rows[0] {
		switch strings.TrimSpace(strings.ToLower(col)) {
		case "no. pesanan":
			colIdx["order_number"] = i
		case "waktu pesanan dibuat":
			colIdx["created_date"] = i
		case "nama produk":
			colIdx["product_name"] = i
		case "jumlah":
			colIdx["quantity"] = i
		}
	}
	if len(colIdx) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required columns in Excel"})
		return
	}

	results := []map[string]interface{}{}
	productNames := []string{}

	type RowData struct {
		RowIdx      int
		OrderNumber string
		CreatedDate time.Time
		ProductName string
		Quantity    int
	}

	type OrderData struct {
		OrderNumber string
		CreatedDate time.Time
		Products    []RowData
	}
	orderMap := make(map[string]*OrderData)

	// First pass: collect all product names and row data
	for rowIdx, row := range rows[1:] {
		orderNumber := row[colIdx["order_number"]]
		createdDateRaw := row[colIdx["created_date"]]
		productName := strings.ToUpper(row[colIdx["product_name"]])
		quantityStr := row[colIdx["quantity"]]
		quantity, err := strconv.Atoi(quantityStr)
		if err != nil {
			results = append(results, map[string]interface{}{"row": rowIdx + 2, "error": "Invalid quantity"})
			continue
		}
		// Parse date (handle both date and time formats)
		createdDate, err := time.Parse("2006-01-02 15:04", createdDateRaw)
		if err != nil {
			createdDate, err = time.Parse("2006-01-02", createdDateRaw)
		}
		if err != nil {
			log.Printf("Failed to parse date '%s': %v", createdDateRaw, err)
			results = append(results, map[string]interface{}{"row": rowIdx + 2, "error": "Invalid created date"})
			continue
		}

		productNames = append(productNames, productName)

		// Group by order number
		if order, exists := orderMap[orderNumber]; exists {
			order.Products = append(order.Products, RowData{
				RowIdx:      rowIdx + 2,
				OrderNumber: orderNumber,
				CreatedDate: createdDate,
				ProductName: productName,
				Quantity:    quantity,
			})
		} else {
			orderMap[orderNumber] = &OrderData{
				OrderNumber: orderNumber,
				CreatedDate: createdDate,
				Products: []RowData{{
					RowIdx:      rowIdx + 2,
					OrderNumber: orderNumber,
					CreatedDate: createdDate,
					ProductName: productName,
					Quantity:    quantity,
				}},
			}
		}
	}

	// Get the first order's created date (all orders should have the same date)
	if len(orderMap) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid orders found in Excel"})
		return
	}

	// Get the first order to get the created date
	var firstOrder *OrderData
	for _, order := range orderMap {
		firstOrder = order
		break
	}

	// Soft delete existing transactions for this date before processing any orders
	username, _ := util.GetUsernameFromContext(c)
	now := time.Now().UTC()
	_, err = h.db.Exec(c, `
		UPDATE online_transactions 
		SET deleted_at = $1, deleted_by = $2
		WHERE DATE(created_date) = DATE($3)
		AND deleted_at IS NULL
	`, now, username, firstOrder.CreatedDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to soft delete existing transactions: " + err.Error()})
		return
	}

	// Query all products at once
	productMap := make(map[string]products.Product)

	productRows, err := h.db.Query(c, `
		SELECT id, name, cost_price, gross_profit_percentage, shopee_category, shopee_name 
		FROM products 
		WHERE shopee_name = ANY($1::text[]) 
		AND deleted_at IS NULL
	`, productNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query products"})
		return
	}
	defer productRows.Close()

	for productRows.Next() {
		var product products.Product
		err := productRows.Scan(
			&product.ID,
			&product.Name,
			&product.CostPrice,
			&product.GrossProfitPercentage,
			&product.ShopeeCategory,
			&product.ShopeeName,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan product data"})
			return
		}
		productMap[product.ShopeeName] = product
	}

	// Process each order
	for _, order := range orderMap {
		// Calculate totals for all products in the order
		var transactionProducts []TransactionProduct
		var hasError bool
		var errorRow int
		var errorMsg string

		// Start a transaction for error logging
		logTx, err := h.db.Begin(c)
		if err != nil {
			results = append(results, map[string]interface{}{"row": order.Products[0].RowIdx, "error": "Failed to begin error logging transaction"})
			continue
		}
		defer logTx.Rollback(c)

		for _, rowData := range order.Products {
			product, exists := productMap[rowData.ProductName]
			if !exists {
				hasError = true
				errorRow = rowData.RowIdx
				errorMsg = "Product not found: " + rowData.ProductName

				// Log the error to database
				_, err = logTx.Exec(c, `
					INSERT INTO product_error_logs (
						order_number, created_date, product_name, error_message, created_by
					) VALUES ($1, $2, $3, $4, $5)
				`, order.OrderNumber, order.CreatedDate, rowData.ProductName, errorMsg, username)
				if err != nil {
					results = append(results, map[string]interface{}{"row": rowData.RowIdx, "error": "Failed to log product error: " + err.Error()})
					logTx.Rollback(c)
					break
				}
				break
			}

			productSalePrice, productFee := products.CalculateShopeePricing(product.CostPrice, product.GrossProfitPercentage, product.ShopeeCategory)
			transactionProducts = append(transactionProducts, TransactionProduct{
				ProductName: product.Name,
				CostPrice:   product.CostPrice,
				SalePrice:   productSalePrice,
				Quantity:    rowData.Quantity,
				FeeAmount:   productFee,
			})
		}

		// Commit the error logging transaction
		if err := logTx.Commit(c); err != nil {
			results = append(results, map[string]interface{}{"row": order.Products[0].RowIdx, "error": "Failed to commit error logging transaction"})
			continue
		}

		if hasError {
			results = append(results, map[string]interface{}{"row": errorRow, "error": errorMsg})
			continue
		}

		totals := calculateTransactionTotals(transactionProducts)

		tx, err := h.db.Begin(c)
		if err != nil {
			results = append(results, map[string]interface{}{"row": order.Products[0].RowIdx, "error": "Failed to begin transaction"})
			continue
		}
		defer tx.Rollback(c)

		username, _ := util.GetUsernameFromContext(c)
		periodMonth := int(order.CreatedDate.Month())
		periodYear := order.CreatedDate.Year()
		txID := uuid.New()
		now := time.Now().UTC()

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
			"SHOPEE",
			order.OrderNumber,
			order.CreatedDate,
			periodMonth,
			periodYear,
			totals.TotalBaseAmount,
			totals.TotalSaleAmount,
			totals.TotalNetProfit,
			username,
			now,
			totals.TotalFeeAmount,
		)
		if err != nil {
			results = append(results, map[string]interface{}{"row": order.Products[0].RowIdx, "error": "Failed to insert transaction: " + err.Error()})
			tx.Rollback(c)
			continue
		}

		// Insert products using bulk insert
		if len(totals.Products) > 0 {
			productQuery := `
				INSERT INTO online_transaction_products (
					id, online_transaction_id, product_name, cost_price, sale_price, quantity, fee_amount
				)
				VALUES 
			`
			values := []interface{}{}
			valueStrings := []string{}
			paramCount := 1

			for _, product := range totals.Products {
				valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
					paramCount, paramCount+1, paramCount+2, paramCount+3, paramCount+4, paramCount+5, paramCount+6))
				values = append(values,
					uuid.New(),
					txID,
					product.ProductName,
					product.CostPrice,
					product.SalePrice,
					product.Quantity,
					product.FeeAmount,
				)
				paramCount += 7
			}

			productQuery += strings.Join(valueStrings, ",")
			_, err = tx.Exec(c, productQuery, values...)
			if err != nil {
				results = append(results, map[string]interface{}{"row": order.Products[0].RowIdx, "error": "Failed to insert products: " + err.Error()})
				tx.Rollback(c)
				continue
			}
		}

		if err := tx.Commit(c); err != nil {
			results = append(results, map[string]interface{}{"row": order.Products[0].RowIdx, "error": "Failed to commit transaction"})
			continue
		}

		results = append(results, map[string]interface{}{"row": order.Products[0].RowIdx, "status": "success", "order_number": order.OrderNumber})
	}
	c.JSON(http.StatusOK, gin.H{"results": results})
}
