package service

import "github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/repository"

type ProductCategoryRules struct {
	repository.ProductCategoryRules
}

func NewProductCategoryRules(
	productCategoryRules repository.ProductCategoryRules,
) *ProductCategoryRules {
	return &ProductCategoryRules{
		ProductCategoryRules: productCategoryRules,
	}
}
