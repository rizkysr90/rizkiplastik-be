package service

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqUpdateRule struct {
	*model.UpdateSizeUnitRulesRequest
}

func (r *reqUpdateRule) sanitize() {
	r.UpdateSizeUnitRulesRequest.SizeUnitID = strings.TrimSpace(r.UpdateSizeUnitRulesRequest.SizeUnitID)
	r.UpdateSizeUnitRulesRequest.ProductCategoryID = strings.TrimSpace(r.UpdateSizeUnitRulesRequest.ProductCategoryID)
	r.UpdateSizeUnitRulesRequest.RuleID = strings.TrimSpace(r.UpdateSizeUnitRulesRequest.RuleID)
}
func (req *reqUpdateRule) validate(ctx context.Context) error {
	fieldValidations := []httperror.FieldValidation{}
	if err := common.ValidateUUIDFormat(req.SizeUnitID); err != nil {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fieldSizeUnitID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateUUIDFormat(req.ProductCategoryID); err != nil {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fieldProductCategoryID,
			Message: err.Error(),
		})
	}
	if len(fieldValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidations)
	}
	return nil
}
func (s *ProductSizeUnitRulesService) UpdateRule(
	ctx context.Context,
	request *model.UpdateSizeUnitRulesRequest,
) error {
	userID := ctx.Value("userID").(string)
	input := &reqUpdateRule{
		UpdateSizeUnitRulesRequest: request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return err
	}
	updatedData := &repository.ProductSizeUnitRulesData{
		RuleID:            input.RuleID,
		SizeUnitID:        input.SizeUnitID,
		ProductCategoryID: input.ProductCategoryID,
		IsDefault:         input.IsDefault,
		UpdatedBy:         userID,
	}
	if err := s.productSizeUnitRulesRepository.UpdateTransaction(ctx, updatedData); err != nil {
		return handleRepositoryError(ctx, err)
	}
	return nil
}
