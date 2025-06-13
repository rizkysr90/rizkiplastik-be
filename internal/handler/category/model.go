package category

import (
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

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

type GetListCategoryRequest struct {
	util.PaginationData `json:"pagination"`
	CategoryName        string `json:"category_name"`
	CategoryCode        string `json:"category_code"`
	IsActive            string `json:"is_active"`
}

type GetListCategoryResponse struct {
	util.PaginationData `json:"pagination"`
	Data                []CategoryBaseModel `json:"data"`
}

type GetByCategoryIDRequest struct {
	CategoryID string `json:"category_id"`
}

type GetByCategoryIDResponse struct {
	CategoryDetailModel `json:"data"`
}
