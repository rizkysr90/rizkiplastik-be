package products_old

import "math"

// CalculateShopeePricing calculates the Shopee sale price and fee based on the formula
func CalculateShopeePricing(
	costPrice, grossProfitPercentage, shopeeFreeDeliveryFee float32,
	category string,
) (salePrice, fee float32) {
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
	deliveryFreeServiceFee := shopeeFreeDeliveryFee / 100
	serviceFee := deliveryFreeServiceFee * grossProfitPriceTotal
	if serviceFee > 20000 {
		serviceFee = 20000
	}

	// 3. Calculate total shopee fee
	fee = adminFee + serviceFee

	// 4. Calculate sale price
	salePrice = fee + grossProfitPriceTotal

	// 5. Round to 2 decimal places
	fee = float32(math.Round(float64(fee)*100) / 100)
	salePrice = float32(math.Round(float64(salePrice)*100) / 100)

	return salePrice, fee
}
