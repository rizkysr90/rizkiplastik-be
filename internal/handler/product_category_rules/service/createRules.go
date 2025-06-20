package service

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqCreateRules struct {
	*model.CreateRulesRequest
}

func (req *reqCreateRules) sanitize() {
	req.PackagingTypeID = strings.TrimSpace(req.PackagingTypeID)
	req.ProductCategoryID = strings.TrimSpace(req.ProductCategoryID)
}
func (req *reqCreateRules) validate(ctx context.Context) error {
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
	if len(fieldValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidations)
	}
	return nil
}
func (s *ProductCategoryRules) CreateRules(
	ctx context.Context,
	request *model.CreateRulesRequest,
) error {
	input := &reqCreateRules{
		CreateRulesRequest: request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return err
	}
	insertedData := &repository.ProductCategoryRulesData{
		RuleID:          uuid.NewString(),
		PackagingTypeID: input.PackagingTypeID,
		CategoryID:      input.ProductCategoryID,
		CreatedBy:       "SYSTEM",
		IsDefault:       false,
		UpdatedBy:       "SYSTEM",
	}
	if input.IsDefault != nil {
		insertedData.IsDefault = *input.IsDefault
	}

	if err := s.ProductCategoryRules.InsertTransaction(ctx, insertedData); err != nil {
		return handleRepositoryError(ctx, err)
	}
	if err := s.ProductCategoryRules.UpdateTransaction(ctx, insertedData); err != nil {
		return handleRepositoryError(ctx, err)
	}
	return nil
}
