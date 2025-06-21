package service

import (
	"context"
	"errors"

	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/repository/pg"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type ProductCategoryRules struct {
	repository.ProductCategoryRules
}

const (
	fieldPackagingTypeID   = "packaging_type_id"
	fieldProductCategoryID = "product_category_id"
	fieldRuleID            = "rule_id"
	fieldStatus            = "status"
)

func NewProductCategoryRules(
	productCategoryRules repository.ProductCategoryRules,
) *ProductCategoryRules {
	return &ProductCategoryRules{
		ProductCategoryRules: productCategoryRules,
	}
}

func handleRepositoryError(ctx context.Context, err error) error {
	if errors.Is(err, pg.ErrCategoryNotFound) {
		return httperror.NewDataNotFound(ctx, httperror.WithMessage(err.Error()))
	}
	if errors.Is(err, pg.ErrPackagingTypeNotFound) {
		return httperror.NewDataNotFound(ctx, httperror.WithMessage(err.Error()))
	}
	if errors.Is(err, pg.ErrRuleAlreadyExists) {
		return httperror.NewBadRequest(ctx, httperror.WithMessage(err.Error()))
	}
	if errors.Is(err, pg.ErrUniqueViolation) {
		return httperror.NewBadRequest(ctx, httperror.WithMessage(err.Error()))
	}
	if errors.Is(err, pg.ErrRuleNotFound) {
		return httperror.NewDataNotFound(ctx, httperror.WithMessage(err.Error()))
	}
	return httperror.NewInternalServer(ctx, httperror.WithMessage(err.Error()))
}
