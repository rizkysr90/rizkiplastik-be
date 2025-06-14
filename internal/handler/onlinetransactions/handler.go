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
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
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
		SELECT id, name, cost_price, gross_profit_percentage, shopee_category, shopee_fee_free_delivery_fee
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
			&temp.ShopeeFreeDeliveryFee,
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
		productSalePrice, productFee := products.CalculateShopeePricing(
			product.CostPrice,
			product.GrossProfitPercentage,
			product.ShopeeFreeDeliveryFee,
			product.ShopeeCategory,
		)
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

// handleExcelFileUpload handles the Excel file upload process and returns the excelize.File
func handleExcelFileUpload(c *gin.Context) (*excelize.File, error) {
	file, err := c.FormFile("file")
	if err != nil {
		return nil, fmt.Errorf("file is required: %w", err)
	}

	f, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	xl, err := excelize.OpenReader(f)
	if err != nil {
		return nil, fmt.Errorf("failed to read Excel file: %w", err)
	}

	return xl, nil
}

// ExcelColumnIndexes represents the column indexes in the Excel file
type ExcelColumnIndexes struct {
	OrderNumber      int
	CreatedDate      int
	ProductName      int
	Quantity         int
	ShopeeVarianName int
	OrderStatus      int
}

// RowData represents a single row of data from the Excel file
type RowData struct {
	RowIdx           int
	OrderNumber      string
	CreatedDate      time.Time
	ProductName      string
	ShopeeVarianName string
	Quantity         int
}

// OrderData represents a group of products for a single order
type OrderData struct {
	OrderNumber string
	CreatedDate time.Time
	Products    []RowData
}

// findColumnIndexes finds the indexes of required columns in the Excel file
func findColumnIndexes(headerRow []string) (ExcelColumnIndexes, error) {
	colIdx := ExcelColumnIndexes{}
	for i, col := range headerRow {
		switch strings.TrimSpace(strings.ToLower(col)) {
		case "no. pesanan":
			colIdx.OrderNumber = i
		case "waktu pesanan dibuat":
			colIdx.CreatedDate = i
		case "nama produk":
			colIdx.ProductName = i
		case "jumlah":
			colIdx.Quantity = i
		case "nama variasi":
			colIdx.ShopeeVarianName = i
		case "status pesanan":
			colIdx.OrderStatus = i
		}
	}

	// Validate required columns
	if colIdx.OrderNumber == 0 && colIdx.CreatedDate == 0 && colIdx.ProductName == 0 && colIdx.Quantity == 0 {
		return colIdx, fmt.Errorf("missing required columns in Excel")
	}

	return colIdx, nil
}

// parseExcelRow parses a single row from the Excel file
func parseExcelRow(row []string, colIdx ExcelColumnIndexes, rowIdx int) (RowData, error) {
	orderNumber := row[colIdx.OrderNumber]
	createdDateRaw := row[colIdx.CreatedDate]
	productName := strings.ToUpper(row[colIdx.ProductName])
	quantityStr := row[colIdx.Quantity]
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		return RowData{}, fmt.Errorf("invalid quantity at row %d", rowIdx+2)
	}

	shopeeVarianName := ""
	if colIdx.ShopeeVarianName < len(row) {
		shopeeVarianName = strings.ToUpper(strings.TrimSpace(row[colIdx.ShopeeVarianName]))
	}

	// Parse date (handle both date and time formats)
	createdDate, err := time.Parse("2006-01-02 15:04", createdDateRaw)
	if err != nil {
		createdDate, err = time.Parse("2006-01-02", createdDateRaw)
	}
	if err != nil {
		return RowData{}, fmt.Errorf("invalid created date at row %d", rowIdx+2)
	}

	return RowData{
		RowIdx:           rowIdx + 2,
		OrderNumber:      orderNumber,
		CreatedDate:      createdDate,
		ProductName:      productName,
		ShopeeVarianName: shopeeVarianName,
		Quantity:         quantity,
	}, nil
}

// groupOrdersByOrderNumber groups rows by order number
func groupOrdersByOrderNumber(rows []RowData) map[string]*OrderData {
	orderMap := make(map[string]*OrderData)
	for _, row := range rows {
		if order, exists := orderMap[row.OrderNumber]; exists {
			order.Products = append(order.Products, row)
		} else {
			orderMap[row.OrderNumber] = &OrderData{
				OrderNumber: row.OrderNumber,
				CreatedDate: row.CreatedDate,
				Products:    []RowData{row},
			}
		}
	}
	return orderMap
}

// queryProducts retrieves all products from the database based on product names
func (h *Handler) queryProducts(c *gin.Context, productNames []string) (map[string]products.Product, error) {
	productMap := make(map[string]products.Product)
	productRows, err := h.db.Query(c, `
		SELECT id, name, cost_price, 
		gross_profit_percentage, 
		shopee_category, COALESCE(shopee_varian_name, ''), shopee_name,
		shopee_fee_free_delivery_fee
		FROM products 
		WHERE shopee_name = ANY($1::text[])
		AND deleted_at IS NULL
	`, productNames)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
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
			&product.ShopeeVarianName,
			&product.ShopeeName,
			&product.ShopeeFreeDeliveryFee,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product data: %w", err)
		}
		lookupKey := product.ShopeeName
		if product.ShopeeVarianName != "" {
			lookupKey = product.ShopeeName + " " + product.ShopeeVarianName
		}
		if product.ShopeeName == "COLATTA DARK CHOCOLATE STICK BAKE STABLE REPACK 100GR 250GR | BEKASI" {
			log.Println("HEREEE : ", product.ShopeeVarianName)
			log.Println("LOOKEY HANDLER : ", lookupKey)
		}
		productMap[lookupKey] = product
	}

	return productMap, nil
}

// processOrder processes a single order and returns transaction products and error status
func (h *Handler) processOrder(
	order *OrderData,
	productMap map[string]products.Product,
	tx pgx.Tx,
	c *gin.Context,
) ([]TransactionProduct, bool, error, []map[string]interface{}) {
	var transactionProducts []TransactionProduct
	username, _ := util.GetUsernameFromContext(c)
	hasError := false
	errorProductNotFound := []map[string]interface{}{}
	for _, rowData := range order.Products {
		lookupKey := rowData.ProductName

		if rowData.ShopeeVarianName != "" {
			lookupKey = rowData.ProductName + " " + rowData.ShopeeVarianName
		}

		product, exists := productMap[lookupKey]
		if rowData.ProductName == strings.ToUpper("Colatta Dark Chocolate Stick Bake Stable Repack 100gr 250gr | Bekasi") {
			log.Println("HEREEE raww: ", rowData.ProductName, rowData.ShopeeVarianName)
			log.Println("RESULT :", product.ShopeeName, product.ShopeeVarianName)
			log.Println("LOOK UP KEY :", lookupKey)
		}
		if !exists {
			errorProductNotFound = append(errorProductNotFound, map[string]interface{}{
				"row":          rowData.RowIdx,
				"error":        "Product not found: " + lookupKey,
				"order_number": order.OrderNumber,
			})
			hasError = true
			// Log the error to database
			_, err := tx.Exec(c, `
				INSERT INTO product_error_logs (
					order_number, created_date, product_name, error_message, created_by
				) VALUES ($1, $2, $3, $4, $5)
			`, order.OrderNumber, order.CreatedDate, lookupKey, "Product not found: "+lookupKey, username)
			if err != nil {
				return nil, true, fmt.Errorf("failed to log product error: %w", err), errorProductNotFound
			}
			continue
		}

		grossProfitPercentage := product.GrossProfitPercentage

		productSalePrice, productFee := products.CalculateShopeePricing(
			product.CostPrice, grossProfitPercentage, product.ShopeeFreeDeliveryFee,
			product.ShopeeCategory)
		transactionProducts = append(transactionProducts, TransactionProduct{
			ProductName: product.Name,
			CostPrice:   product.CostPrice,
			SalePrice:   productSalePrice,
			Quantity:    rowData.Quantity,
			FeeAmount:   productFee,
		})
	}

	return transactionProducts, hasError, nil, errorProductNotFound
}

// insertTransaction inserts a transaction and its products into the database
func insertTransaction(order *OrderData, totals TransactionTotals, tx pgx.Tx, c *gin.Context) error {
	username, _ := util.GetUsernameFromContext(c)
	now := time.Now().UTC()
	uuidOnlineTransaction := uuid.New()

	// Insert transaction
	transactionQuery := `
		INSERT INTO online_transactions (
			id, type, order_number, created_date, period_month, period_year,
			total_base_amount, total_sale_amount, total_net_profit, created_by,
			created_at, total_fee_amount
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := tx.Exec(c, transactionQuery,
		uuidOnlineTransaction,
		"SHOPEE",
		order.OrderNumber,
		order.CreatedDate,
		int(order.CreatedDate.Month()),
		order.CreatedDate.Year(),
		totals.TotalBaseAmount,
		totals.TotalSaleAmount,
		totals.TotalNetProfit,
		username,
		now,
		totals.TotalFeeAmount,
	)
	if err != nil {
		return fmt.Errorf("failed to insert transaction: %w", err)
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
				uuidOnlineTransaction,
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
			return fmt.Errorf("failed to insert products: %w", err)
		}
	}

	return nil
}

// processTransaction handles the entire transaction process including soft deletes and product processing
func (h *Handler) processTransaction(orderMap map[string]*OrderData, productMap map[string]products.Product, c *gin.Context) {
	username, _ := util.GetUsernameFromContext(c)
	now := time.Now().UTC()

	// Start transaction
	tx, err := h.db.Begin(c)
	if err != nil {
		errMessage := fmt.Sprintf("failed to begin transaction: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}
	defer tx.Rollback(c)

	// Get the first order to get the created date
	var firstOrder *OrderData
	for _, order := range orderMap {
		firstOrder = order
		break
	}

	// Soft delete existing transactions and error logs
	_, err = tx.Exec(c, `
		UPDATE online_transactions 
		SET deleted_at = $1, deleted_by = $2
		WHERE DATE(created_date) = DATE($3)
		AND deleted_at IS NULL
	`, now, username, firstOrder.CreatedDate)
	if err != nil {
		errMessage := fmt.Sprintf("failed to soft delete existing transactions: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}

	_, err = tx.Exec(c, `
		UPDATE product_error_logs 
		SET deleted_at = $1
		WHERE DATE(created_date) = DATE($2)
		AND deleted_at IS NULL
	`, now, firstOrder.CreatedDate)
	if err != nil {
		errMessage := fmt.Sprintf("failed to soft delete existing product error logs: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}

	// Process each order
	errorResultsProductNotFound := []map[string]interface{}{}
	for _, order := range orderMap {
		transactionProducts, hasError, err, errorProductNotFound := h.processOrder(order, productMap, tx, c)
		if err != nil {
			errMessage := fmt.Sprintf("failed to process order: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
			return
		}

		if hasError {
			errorResultsProductNotFound = append(errorResultsProductNotFound, errorProductNotFound...)
			continue
		}

		totals := calculateTransactionTotals(transactionProducts)
		err = insertTransaction(order, totals, tx, c)
		if err != nil {
			errMessage := fmt.Sprintf("failed to insert transaction: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
			return
		}
	}
	if len(errorResultsProductNotFound) > 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found", "results": errorResultsProductNotFound})
		return
	}
	// Commit transaction
	if err := tx.Commit(c); err != nil {
		errMessage := fmt.Sprintf("failed to commit transaction: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": errMessage})
		return
	}
	c.JSON(http.StatusOK, gin.H{})
}

// AutoInputOnlineTransactions handles Excel upload and auto input
func (h *Handler) AutoInputOnlineTransactions(c *gin.Context) {
	xl, err := handleExcelFileUpload(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sheetName := xl.GetSheetName(0)
	rows, err := xl.GetRows(sheetName)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read rows from Excel"})
		return
	}

	// Find column indexes
	colIdx, err := findColumnIndexes(rows[0])
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse all rows
	var allRows []RowData
	productNames := []string{}
	for rowIdx, row := range rows[1:] {
		if strings.ToLower(row[colIdx.OrderStatus]) == "batal" {
			continue
		}
		rowData, err := parseExcelRow(row, colIdx, rowIdx)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		allRows = append(allRows, rowData)
		productNames = append(productNames, rowData.ProductName)
	}

	// Group orders
	orderMap := groupOrdersByOrderNumber(allRows)
	if len(orderMap) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No valid orders found in Excel"})
		return
	}

	// Query all products
	productMap, err := h.queryProducts(c, productNames)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Process all transactions
	h.processTransaction(orderMap, productMap, c)
}
