package service

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/handler/product_category_rules/model"
)

func (s *ProductCategoryRules) UpdateStatusRules(
	ctx context.Context,
	request *model.UpdateRulesStatusRequest,
) error {
	request.RuleID = strings.TrimSpace(request.RuleID)
	if err := s.ProductCategoryRules.UpdateStatusRule(
		ctx, request.RuleID, request.Status); err != nil {
		return handleRepositoryError(ctx, err)
	}
	return nil
}
