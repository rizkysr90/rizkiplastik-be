package model

type SizeUnit struct {
	SizeUnitName string `json:"size_unit_name"`
	SizeUnitCode string `json:"size_unit_code"`
	SizeUnitType string `json:"size_unit_type"`
}

// Base Model
type SizeUnitRules struct {
	RuleID            string   `json:"rule_id"`
	ProductCategoryID string   `json:"product_category_id"`
	SizeUnitID        string   `json:"size_unit_id"`
	SizeUnit          SizeUnit `json:"size_unit"`
	IsDefault         bool     `json:"is_default"`
	IsActive          bool     `json:"is_active"`
}
