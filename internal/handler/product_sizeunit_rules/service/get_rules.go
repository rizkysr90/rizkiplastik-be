package service

import (
	"context"
	"html"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqGetRules struct {
	*model.GetListSizeUnitRulesRequest
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
func (s *ProductSizeUnitRulesService) GetRules(
	ctx context.Context,
	request *model.GetListSizeUnitRulesRequest,
) (*model.GetListSizeUnitRulesResponse, error) {
	input := &reqGetRules{
		GetListSizeUnitRulesRequest: request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return nil, err
	}
	filter := repository.ProductSizeUnitRulesFilter{
		CategoryID: input.ProductCategoryID,
		Status:     input.Status,
	}
	rules, err := s.productSizeUnitRulesRepository.FindSizeUnitRulesByCategoryID(ctx, filter)
	if err != nil {
		return nil, handleRepositoryError(ctx, err)
	}
	response := model.GetListSizeUnitRulesResponse{
		Data: []model.SizeUnitRules{},
	}
	for _, rule := range rules {
		response.Data = append(response.Data, model.SizeUnitRules{
			RuleID:            rule.RuleID,
			ProductCategoryID: rule.ProductCategoryID,
			SizeUnitID:        rule.SizeUnitID,
			SizeUnit: model.SizeUnit{
				SizeUnitType: rule.SizeUnitType,
				SizeUnitCode: rule.SizeUnitCode,
				SizeUnitName: html.EscapeString(rule.SizeUnitName),
			},
			IsDefault: rule.IsDefault,
			IsActive:  rule.IsActive,
		})
	}
	return &response, nil
}
