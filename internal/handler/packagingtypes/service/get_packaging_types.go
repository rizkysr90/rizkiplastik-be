package service

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqGetPackagingTypes struct {
	*model.RequestGetPackagingTypes
}

func (req *reqGetPackagingTypes) sanitize() {
	req.Name = strings.TrimSpace(strings.ToUpper(req.Name))
	req.Code = strings.TrimSpace(strings.ToUpper(req.Code))
	req.IsActive = strings.TrimSpace(strings.ToUpper(req.IsActive))
}
func (req *reqGetPackagingTypes) validateFields(ctx context.Context) error {
	fieldValidations := []httperror.FieldValidation{}
	if req.Name != "" {
		fieldValidations = append(fieldValidations, validateCategoryTypeName(req.Name)...)
	}
	if req.Code != "" {
		fieldValidations = append(fieldValidations, validateCategoryTypeCode(req.Code)...)
	}
	if req.IsActive != "" {
		if err := common.ValidateEquals(req.IsActive, []string{
			constants.IsActiveTrue, constants.IsActiveFalse}); err != nil {
			fieldValidations = append(fieldValidations,
				httperror.NewFieldValidation(fieldPackagingIsActive, err.Error()))
		}
	}
	if len(fieldValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidations)
	}
	return nil
}

func (s *PackagingType) GetPackagingTypes(ctx context.Context,
	request *model.RequestGetPackagingTypes) (*model.ResponseGetPackagingTypes, error) {
	input := &reqGetPackagingTypes{
		RequestGetPackagingTypes: request,
	}
	input.sanitize()
	if err := input.validateFields(ctx); err != nil {
		return nil, err
	}
	filter := &repository.PackagingTypeFilter{
		Name:     input.Name,
		Code:     input.Code,
		IsActive: input.IsActive,
		Limit:    input.PageSize,
		Offset:   input.GetOffset(),
	}
	packagingTypes, totalCount, err := s.packagingTypeRepo.FindPaginatedPackagingTypes(ctx, filter)
	if err != nil {
		return nil, err
	}
	input.SetTotalPagesAndTotalElement(totalCount)
	paginationData := &util.PaginationData{
		PageNumber:    input.PageNumber,
		PageSize:      input.PageSize,
		TotalElements: input.TotalElements,
		TotalPages:    input.TotalPages,
	}
	response := make([]model.SimplePackagingType, len(packagingTypes))
	for i, packagingType := range packagingTypes {
		response[i] = model.SimplePackagingType{
			PackagingID:   packagingType.ID,
			PackagingName: packagingType.Name,
			PackagingCode: packagingType.Code,
			IsActive:      packagingType.IsActive,
			CreatedAt:     packagingType.CreatedAt,
			UpdatedAt:     packagingType.UpdatedAt,
		}
	}
	return &model.ResponseGetPackagingTypes{
		PaginationData: paginationData,
		Data:           response,
	}, nil
}
