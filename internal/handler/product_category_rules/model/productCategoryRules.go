package model

type CreateRulesRequest struct {
	ProductCategoryID string `json:"product_category_id"`
	PackagingTypeID   string `json:"packaging_type_id"`
	IsDefault         *bool  `json:"is_default"`
}

type UpdateRulesRequest struct {
	RuleID            string
	ProductCategoryID string
	PackagingTypeID   string `json:"packaging_type_id"`
	IsDefault         bool   `json:"is_default"`
}

type GetListRulesRequest struct {
	ProductCategoryID string `json:"product_category_id"`
	Status            string `json:"status"`
}

type UpdateRulesStatusRequest struct {
	ProductCategoryID string `json:"product_category_id"`
	RuleID            string `json:"rule_id"`
	Status            string `json:"status"`
}

type GetListRulesResponse struct {
	Data []Rules `json:"data"`
}
