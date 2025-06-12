package category

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

type reqUpdateCategory struct {
	*UpdateCategoryRequest
}

func (req *reqUpdateCategory) sanitize() {
	req.CategoryID = strings.TrimSpace(req.CategoryID)
	req.CategoryName = strings.TrimSpace(strings.ToUpper(req.CategoryName))
	if req.CategoryDescription != nil {
		*req.CategoryDescription = strings.TrimSpace(strings.ToUpper(*req.CategoryDescription))
	}
}

func (req *reqUpdateCategory) validate() error {
	if err := common.ValidateUUIDFormat(req.CategoryID); err != nil {
		return &util.ServiceError{
			HTTPCode: 400,
			Message:  err.Error(),
		}
	}
	if err := validateCategoryName(req.CategoryName); err != nil {
		return err
	}
	if req.CategoryDescription != nil {
		if err := validateCategoryDescription(*req.CategoryDescription); err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) UpdateCategory(ctx context.Context, data *UpdateCategoryRequest) error {
	input := &reqUpdateCategory{
		UpdateCategoryRequest: data,
	}
	input.sanitize()
	if err := input.validate(); err != nil {
		return err
	}
	updatedCategory := &repository.CategoryData{
		ID:        input.CategoryID,
		Name:      input.CategoryName,
		IsActive:  input.IsActive,
		UpdatedBy: "SYSTEM",
	}
	if input.CategoryDescription != nil {
		updatedCategory.Description = *input.CategoryDescription
	}
	if err := s.categoryRepo.Update(ctx, updatedCategory); err != nil {
		return util.ConvertRepositoryError(err)
	}
	return nil
}
