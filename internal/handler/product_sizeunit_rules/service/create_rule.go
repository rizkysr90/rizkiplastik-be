package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqCreateSizeUnitRules struct {
	*model.CreateSizeUnitRulesRequest
}

func (req *reqCreateSizeUnitRules) sanitize() {
	req.ProductCategoryID = strings.TrimSpace(
		req.ProductCategoryID)
	req.SizeUnitID = strings.TrimSpace(
		req.SizeUnitID)
}
func (req *reqCreateSizeUnitRules) validate(ctx context.Context) error {
	fieldValidations := []httperror.FieldValidation{}
	if err := common.ValidateUUIDFormat(req.ProductCategoryID); err != nil {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fieldProductCategoryID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateUUIDFormat(req.SizeUnitID); err != nil {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fieldSizeUnitID,
			Message: err.Error(),
		})
	}
	if len(fieldValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx,
			fieldValidations)
	}
	return nil
}
func (s *ProductSizeUnitRulesService) CreateRule(
	ctx context.Context,
	request *model.CreateSizeUnitRulesRequest,
) error {
	userID := ctx.Value("userID").(string)
	input := &reqCreateSizeUnitRules{
		CreateSizeUnitRulesRequest: request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return err
	}
	data := &repository.ProductSizeUnitRulesData{
		RuleID:            uuid.NewString(),
		ProductCategoryID: input.ProductCategoryID,
		SizeUnitID:        input.SizeUnitID,
		CreatedBy:         userID,
		IsDefault:         false,
		UpdatedBy:         userID,
	}
	if input.IsDefault != nil {
		data.IsDefault = *input.IsDefault
	}
	if err := s.productSizeUnitRulesRepository.InsertTransaction(ctx, data); err != nil {
		return handleRepositoryError(ctx, err)
	}
	return nil
}
