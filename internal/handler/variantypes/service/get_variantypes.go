package service

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqGetVarianTypes struct {
	*model.RequestVarianTypePaginated
}

func (r *reqGetVarianTypes) sanitize() {
	r.IsActive = strings.TrimSpace(strings.ToUpper(r.IsActive))
	r.VarianTypeName = strings.TrimSpace(strings.ToUpper(r.VarianTypeName))
}

func (r *reqGetVarianTypes) validate(ctx context.Context) error {
	fieldErrors := []httperror.FieldValidation{}
	if r.VarianTypeName != "" {
		fieldErrors = append(fieldErrors, ValidateVarianTypeName(r.VarianTypeName)...)
	}
	if r.IsActive != "" {
		if err := common.ValidateEquals(r.IsActive, []string{
			constants.IsActiveTrue, constants.IsActiveFalse}); err != nil {
			fieldErrors = append(fieldErrors, httperror.FieldValidation{
				Field:   "is_active",
				Message: err.Error(),
			})
		}
	}

	if len(fieldErrors) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldErrors)
	}
	return nil
}
func (s *VarianTypes) GetVarianTypes(
	ctx context.Context,
	request *model.RequestVarianTypePaginated,
) (*model.ResponseVarianTypePaginated, error) {
	input := &reqGetVarianTypes{
		RequestVarianTypePaginated: request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return nil, err
	}
	filter := repository.VarianTypeFilter{
		Name:     input.VarianTypeName,
		IsActive: input.IsActive,
		Offset:   input.GetOffset(),
		Limit:    input.PageSize,
	}
	variantTypes, totalRows, err := s.repository.FindVarianTypePaginated(ctx, filter)
	if err != nil {
		return nil, err
	}
	input.SetTotalPagesAndTotalElement(totalRows)
	response := &model.ResponseVarianTypePaginated{
		Pagination: &input.PaginationData,
		Data:       make([]model.VarianTypeSimple, 0),
	}
	for _, variantType := range variantTypes {
		response.Data = append(response.Data, model.VarianTypeSimple{
			VarianTypeID:   variantType.ID,
			VarianTypeName: variantType.Name,
			IsActive:       variantType.IsActive,
			CreatedAt:      variantType.CreatedAt,
			UpdatedAt:      variantType.UpdatedAt,
		})
	}
	return response, nil
}
