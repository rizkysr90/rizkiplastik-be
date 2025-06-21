package service

import (
	"context"
	"html"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqGetRules struct {
	*model.GetListRulesRequest
}

func (req *reqGetRules) sanitize() {
	req.ProductCategoryID = strings.TrimSpace(req.ProductCategoryID)
	req.Status = strings.TrimSpace(req.Status)
}
func (req *reqGetRules) validate(ctx context.Context) error {
	fieldValidations := []httperror.FieldValidation{}
	if err := common.ValidateUUIDFormat(req.ProductCategoryID); err != nil {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fieldProductCategoryID,
			Message: err.Error(),
		})
	}
	if req.Status != "" {
		if err := common.ValidateEquals(
			req.Status,
			[]string{constants.IsActiveTrue, constants.IsActiveFalse}); err != nil {
			fieldValidations = append(fieldValidations, httperror.FieldValidation{
				Field:   fieldStatus,
				Message: err.Error(),
			})
		}
	}
	if len(fieldValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidations)
	}
	return nil
}
func (s *ProductCategoryRules) GetRules(
	ctx context.Context,
	request *model.GetListRulesRequest,
) (*model.GetListRulesResponse, error) {
	input := &reqGetRules{
		GetListRulesRequest: request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return nil, err
	}
	filter := repository.ProductCategoryRulesFilter{
		CategoryID: input.ProductCategoryID,
		Status:     input.Status,
	}
	rules, err := s.ProductCategoryRules.FindRuleByCategoryID(ctx, filter)
	if err != nil {
		return nil, handleRepositoryError(ctx, err)
	}
	response := model.GetListRulesResponse{
		Data: []model.Rules{},
	}
	for _, rule := range rules {
		response.Data = append(response.Data, model.Rules{
			RuleID:            rule.RuleID,
			ProductCategoryID: rule.CategoryID,
			PackagingTypeID:   rule.PackagingTypeID,
			PackagingType: model.PackagingType{
				PackagingCode: rule.PackagingCode,
				PackagingName: html.EscapeString(rule.PackagingName),
			},
			IsDefault: rule.IsDefault,
			IsActive:  rule.IsActive,
		})
	}
	return &response, nil
}
