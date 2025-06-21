package model

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
