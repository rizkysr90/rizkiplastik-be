package products

import "math"

// calculateShopeePricing calculates the Shopee sale price and fee based on the formula
func calculateShopeePricing(costPrice, grossProfitPercentage float32, category string) (salePrice, fee float32) {
	// 1. Find gross_profit_price_total = ROUND(cost_price * (gross_profit_percentage / 100) + cost_price)
	grossProfitPriceTotal := float32(math.Round(float64(costPrice*(grossProfitPercentage/100) + costPrice)))

	// 2. Find shopee fee SUM(admin fee + service fee)
	var adminFee float32

	// a. Calculate admin fee based on category
	switch category {
	case "A":
		adminFee = 0.1 * grossProfitPriceTotal
	case "B":
		adminFee = 0.075 * grossProfitPriceTotal
	case "C":
		adminFee = 0.0575 * grossProfitPriceTotal
	default: // D and E
		adminFee = 0.0575 * grossProfitPriceTotal // Using C rate as fallback
	}

	// b. Calculate service fee with max 10,000
	serviceFee := 0.04 * grossProfitPriceTotal
	if serviceFee > 10000 {
		serviceFee = 10000
	}

	// 3. Calculate total shopee fee
	fee = adminFee + serviceFee

	// 4. Calculate sale price
	salePrice = fee + grossProfitPriceTotal

	return salePrice, fee
}
