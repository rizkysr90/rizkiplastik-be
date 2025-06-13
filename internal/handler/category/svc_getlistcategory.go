package category

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

const (
	IsActiveTrue  = "TRUE"
	IsActiveFalse = "FALSE"
	IsActiveAll   = "ALL"
)

type reqGetListCategory struct {
	*GetListCategoryRequest
}

func (req *reqGetListCategory) sanitize() {
	req.CategoryName = strings.TrimSpace(strings.ToUpper(req.CategoryName))
	req.CategoryCode = strings.TrimSpace(strings.ToUpper(req.CategoryCode))
	req.IsActive = strings.TrimSpace(strings.ToUpper(req.IsActive))
}
func (req *reqGetListCategory) validate(ctx context.Context) error {
	fieldValidationErrors := []httperror.FieldValidation{}
	if err := common.ValidateMaxLengthStr(req.CategoryName, 50); err != nil {
		fieldValidationErrors = append(fieldValidationErrors,
			httperror.NewFieldValidation(
				fieldCategoryName, err.Error()))
	}
	if err := common.ValidateMaxLengthStr(req.CategoryCode, 3); err != nil {
		fieldValidationErrors = append(fieldValidationErrors,
			httperror.NewFieldValidation(
				fieldCategoryCode, err.Error()))
	}

	if req.IsActive != IsActiveTrue &&
		req.IsActive != IsActiveFalse &&
		req.IsActive != IsActiveAll {
		fieldValidationErrors = append(fieldValidationErrors,
			httperror.NewFieldValidation(
				fieldIsActive, "invalid is_active value"))
	}
	if len(fieldValidationErrors) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidationErrors)
	}
	return nil
}

func (s *Service) GetListCategory(ctx context.Context,
	request *GetListCategoryRequest) (*GetListCategoryResponse, error) {
	input := &reqGetListCategory{
		GetListCategoryRequest: request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return nil, err
	}
	filter := &repository.CategoryDataFilter{
		CategoryName: input.CategoryName,
		CategoryCode: input.CategoryCode,
		IsActive:     input.IsActive,
		PageSize:     input.PageSize,
		Offset:       input.GetOffset(),
	}
	categories, totalCount, err := s.categoryRepo.GetList(ctx, filter)
	if err != nil {
		return nil, httperror.NewInternalServer(
			ctx,
			httperror.WithMessage("failed to get list category : "+err.Error()),
		)
	}
	categoryBaseModels := make([]CategoryBaseModel, 0)
	for _, category := range categories {
		categoryBaseModels = append(categoryBaseModels, CategoryBaseModel{
			CategoryID:   category.ID,
			CategoryName: category.Name,
			CategoryCode: category.Code,
			IsActive:     category.IsActive,
			CreatedAt:    category.CreatedAt,
			UpdatedAt:    category.UpdatedAt,
		})
	}
	input.PaginationData.SetTotalPagesAndTotalElement(totalCount)
	response := &GetListCategoryResponse{
		Data:           categoryBaseModels,
		PaginationData: input.PaginationData,
	}
	return response, nil
}
