package category

type CreateCategoryRequest struct {
	Name        string  `json:"category_name"`
	Code        string  `json:"category_code"`
	Description *string `json:"category_description,omitempty"`
}

type UpdateCategoryRequest struct {
	CategoryID          string  `json:"category_id"`
	CategoryName        string  `json:"category_name"`
	CategoryDescription *string `json:"category_description,omitempty"`
	IsActive            bool    `json:"is_active"`
}
