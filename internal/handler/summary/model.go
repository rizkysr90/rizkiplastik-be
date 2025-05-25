package summary

type RequestSummary struct {
	StartDate string `json:"start_date" binding:"required"` // ISO 8601 format
	EndDate   string `json:"end_date" binding:"required"`   // ISO 8601 format
}

type DailyProfit struct {
	Date      string  `json:"date"`
	NetProfit float64 `json:"net_profit"`
}

type ResponseSummary struct {
	TotalNetProfit float64       `json:"total_net_profit"`
	DailyProfits   []DailyProfit `json:"daily_profits"`
}
