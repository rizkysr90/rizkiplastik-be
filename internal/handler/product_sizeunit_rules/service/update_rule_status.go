package service

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/model"
)

func (s *ProductSizeUnitRulesService) UpdateRuleStatus(
	ctx context.Context,
	request *model.UpdateSizeUnitRulesStatusRequest,
) error {
	request.RuleID = strings.TrimSpace(request.RuleID)
	if err := s.productSizeUnitRulesRepository.UpdateStatusRule(
		ctx, request.RuleID, request.Status); err != nil {
		return handleRepositoryError(ctx, err)
	}
	return nil
}
