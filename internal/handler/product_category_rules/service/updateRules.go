package service

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqUpdateRules struct {
	*model.UpdateRulesRequest
}

func (r *reqUpdateRules) sanitize() {
	r.UpdateRulesRequest.PackagingTypeID = strings.TrimSpace(r.UpdateRulesRequest.PackagingTypeID)
	r.UpdateRulesRequest.ProductCategoryID = strings.TrimSpace(r.UpdateRulesRequest.ProductCategoryID)
	r.UpdateRulesRequest.RuleID = strings.TrimSpace(r.UpdateRulesRequest.RuleID)
}
func (req *reqUpdateRules) validate(ctx context.Context) error {
	fieldValidations := []httperror.FieldValidation{}
	if err := common.ValidateUUIDFormat(req.PackagingTypeID); err != nil {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fieldPackagingTypeID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateUUIDFormat(req.ProductCategoryID); err != nil {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fieldProductCategoryID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateUUIDFormat(req.RuleID); err != nil {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fieldRuleID,
			Message: err.Error(),
		})
	}

	if len(fieldValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidations)
	}
	return nil
}
func (s *ProductCategoryRules) UpdateRules(
	ctx context.Context,
	request *model.UpdateRulesRequest,
) error {
	input := &reqUpdateRules{
		UpdateRulesRequest: request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return err
	}
	updatedData := &repository.ProductCategoryRulesData{
		RuleID:          input.RuleID,
		PackagingTypeID: input.PackagingTypeID,
		CategoryID:      input.ProductCategoryID,
		IsDefault:       input.IsDefault,
		UpdatedBy:       "SYSTEM",
	}
	if err := s.ProductCategoryRules.UpdateTransaction(ctx, updatedData); err != nil {
		return handleRepositoryError(ctx, err)
	}
	return nil
}
