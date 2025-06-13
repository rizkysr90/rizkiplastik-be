package category

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
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

func (req *reqUpdateCategory) validate(ctx context.Context) error {
	fieldValidationErrors := []httperror.FieldValidation{}
	if err := common.ValidateUUIDFormat(req.CategoryID); err != nil {
		fieldValidationErrors = append(
			fieldValidationErrors,
			httperror.NewFieldValidation(fieldCategoryID, err.Error()))
	}
	fieldValidationErrors = append(
		fieldValidationErrors,
		validateCategoryName(req.CategoryName)...)
	if req.CategoryDescription != nil {
		fieldValidationErrors = append(
			fieldValidationErrors,
			validateCategoryDescription(*req.CategoryDescription)...)
	}
	if len(fieldValidationErrors) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidationErrors)
	}
	return nil
}
func (s *Service) UpdateCategory(ctx context.Context, data *UpdateCategoryRequest) error {
	input := &reqUpdateCategory{
		UpdateCategoryRequest: data,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
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
		if errors.Is(err, pgx.ErrNoRows) {
			return httperror.NewDataNotFound(
				ctx, httperror.WithMessage("category not found"))
		}
		return httperror.NewInternalServer(
			ctx, httperror.WithMessage("failed to update category : "+err.Error()))
	}
	return nil
}
