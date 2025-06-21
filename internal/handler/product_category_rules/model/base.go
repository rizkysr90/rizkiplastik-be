package model

type PackagingType struct {
	PackagingCode string `json:"packaging_code"`
	PackagingName string `json:"packaging_name"`
}

type Rules struct {
	RuleID            string        `json:"rule_id"`
	ProductCategoryID string        `json:"product_category_id"`
	PackagingTypeID   string        `json:"packaging_type_id"`
	PackagingType     PackagingType `json:"packaging_type"`
	IsDefault         bool          `json:"is_default"`
	IsActive          bool          `json:"is_active"`
}
