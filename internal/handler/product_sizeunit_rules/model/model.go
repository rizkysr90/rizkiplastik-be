package model

type SizeUnit struct {
	SizeUnitID   string `json:"size_unit_id"`
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

type CreateSizeUnitRulesRequest struct {
	ProductCategoryID string // in params
	SizeUnitID        string `json:"size_unit_id"`
	IsDefault         *bool  `json:"is_default"`
}

type UpdateSizeUnitRulesRequest struct {
	ProductCategoryID string // in params
	RuleID            string // in params
	SizeUnitID        string `json:"size_unit_id"`
	IsDefault         bool   `json:"is_default"`
}

type GetListSizeUnitRulesRequest struct {
	ProductCategoryID string // in params
	Status            string // in query
}

type GetListSizeUnitRulesResponse struct {
	Data []SizeUnitRules `json:"data"`
}

type UpdateSizeUnitRulesStatusRequest struct {
	// in params
	RuleID string
	// in body
	Status bool `json:"status"`
}
