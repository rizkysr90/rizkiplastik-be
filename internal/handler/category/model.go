package category

type CreateCategoryRequest struct {
	Name        string  `json:"category_name"`
	Code        string  `json:"category_code"`
	Description *string `json:"category_description,omitempty"`
}
