package service

import (
	"context"
	"errors"

	packagingRulesPg "github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/repository/pg"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository/pg"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

const (
	fieldProductCategoryID = "product_category_id"
	fieldSizeUnitID        = "size_unit_id"
	fieldIsDefault         = "is_default"
	fieldRuleID            = "rule_id"
)

type ProductSizeUnitRulesService struct {
	productSizeUnitRulesRepository repository.ProductSizeUnitRules
}

func NewProductSizeUnitRulesService(
	productSizeUnitRulesRepository repository.ProductSizeUnitRules,
) *ProductSizeUnitRulesService {
	return &ProductSizeUnitRulesService{
		productSizeUnitRulesRepository: productSizeUnitRulesRepository,
	}
}
func handleRepositoryError(ctx context.Context, err error) error {
	if errors.Is(err, packagingRulesPg.ErrCategoryNotFound) {
		return httperror.NewDataNotFound(ctx, httperror.WithMessage(err.Error()))
	}
	if errors.Is(err, pg.ErrRuleSizeUnitAlreadyExists) {
		return httperror.NewBadRequest(ctx, httperror.WithMessage(err.Error()))
	}
	if errors.Is(err, pg.ErrUniqueViolation) {
		return httperror.NewBadRequest(ctx, httperror.WithMessage(err.Error()))
	}
	if errors.Is(err, pg.ErrRuleSizeUnitNotFound) {
		return httperror.NewDataNotFound(ctx, httperror.WithMessage(err.Error()))
	}
	if errors.Is(err, pg.ErrSizeUnitNotFound) {
		return httperror.NewDataNotFound(ctx, httperror.WithMessage(err.Error()))
	}
	return httperror.NewInternalServer(ctx, httperror.WithMessage(err.Error()))
}
