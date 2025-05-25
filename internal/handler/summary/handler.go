package summary

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/middleware"
)

// SummaryHandler handles HTTP requests for products
type SummaryHandler struct {
	db *pgxpool.Pool
}

// NewSummaryHandler creates a new product handler
func NewSummaryHandler(db *pgxpool.Pool) *SummaryHandler {
	return &SummaryHandler{db: db}
}

// RegisterRoutes registers all product-related routes
func (h *SummaryHandler) RegisterRoutes(
	router *gin.Engine,
	authMiddleware *middleware.AuthMiddleware) {
	private := router.Group("api/v1/summaries")
	private.Use(authMiddleware.RequireAuth())
	{
		private.GET("", h.GetSummary) // Anyone can list products

	}
}

// Using the standard library time package
func isValidISO8601Date(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}
func (h *SummaryHandler) GetSummary(c *gin.Context) {
	start_date := c.Query("start_date")
	end_date := c.Query("end_date")
	if start_date == "" || end_date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required"})
		return
	}
	request := RequestSummary{
		StartDate: start_date,
		EndDate:   end_date,
	}
	// Sanitize ISO 8601 FORMAT
	request.EndDate = strings.TrimSpace(request.EndDate)
	request.StartDate = strings.TrimSpace(request.StartDate)

	// Validate ISO 8601 FORMAT
	if !isValidISO8601Date(request.EndDate) || !isValidISO8601Date(request.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ISO 8601 format"})
		return
	}

	// Query for total net profit
	var totalNetProfit float64
	totalQuery := `
		SELECT  
			COALESCE(SUM(total_net_profit), 0) as total_net_profit
		FROM 
			online_transactions
		WHERE 
			created_date BETWEEN $1::DATE AND $2::DATE 
			AND deleted_at IS NULL
	`
	err := h.db.QueryRow(c, totalQuery, request.StartDate, request.EndDate).Scan(&totalNetProfit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Query for daily breakdown
	dailyQuery := `
		SELECT  
			TO_CHAR(created_date, 'YYYY-MM-DD') as date,
			COALESCE(SUM(total_net_profit), 0) as net_profit
		FROM 
			online_transactions
		WHERE 
			created_date BETWEEN $1::DATE AND $2::DATE 
			AND deleted_at IS NULL
		GROUP BY DATE(created_date)
		ORDER BY DATE(created_date)
	`
	rows, err := h.db.Query(c, dailyQuery, request.StartDate, request.EndDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var dailyProfits []DailyProfit
	for rows.Next() {
		var daily DailyProfit
		if err := rows.Scan(&daily.Date, &daily.NetProfit); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() + "in scan daily"})
			return
		}
		dailyProfits = append(dailyProfits, daily)
	}

	response := ResponseSummary{
		TotalNetProfit: totalNetProfit,
		DailyProfits:   dailyProfits,
	}
	c.JSON(http.StatusOK, response)
}
