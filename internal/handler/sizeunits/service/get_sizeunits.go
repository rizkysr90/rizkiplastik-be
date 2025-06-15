package service

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqGetSizeUnits struct {
	*model.RequestGetSizeUnits
}

func (r *reqGetSizeUnits) sanitize() {
	r.SizeUnitCode = strings.TrimSpace(strings.ToUpper(r.SizeUnitCode))
	r.SizeUnitName = strings.TrimSpace(strings.ToUpper(r.SizeUnitName))
	r.SizeUnitType = strings.TrimSpace(strings.ToUpper(r.SizeUnitType))
	r.IsActive = strings.TrimSpace(strings.ToUpper(r.IsActive))
}
func (r *reqGetSizeUnits) validate(ctx context.Context) error {
	fieldsValidations := []httperror.FieldValidation{}
	if r.SizeUnitCode != "" {
		fieldsValidations = append(fieldsValidations, ValidateSizeUnitCode(r.SizeUnitCode)...)
	}
	if r.SizeUnitName != "" {
		fieldsValidations = append(fieldsValidations, ValidateSizeUnitName(r.SizeUnitName)...)
	}
	if r.SizeUnitType != "" {
		fieldsValidations = append(fieldsValidations, ValidateSizeUnitType(r.SizeUnitType)...)
	}
	if r.IsActive != "" {
		if err := common.ValidateEquals(r.IsActive, []string{
			constants.IsActiveFalse,
			constants.IsActiveTrue,
		}); err != nil {
			fieldsValidations = append(fieldsValidations, httperror.FieldValidation{
				Field:   "is_active",
				Message: err.Error(),
			})
		}
	}

	if len(fieldsValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldsValidations)
	}
	return nil
}
func (s *SizeUnits) GetSizeUnits(ctx context.Context, request *model.RequestGetSizeUnits) (
	*model.ResponseGetSizeUnits, error) {
	input := &reqGetSizeUnits{
		RequestGetSizeUnits: request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return nil, err
	}
	filter := repository.SizeUnitFilter{
		SizeUnitName: input.SizeUnitName,
		SizeUnitCode: input.SizeUnitCode,
		SizeUnitType: input.SizeUnitType,
		IsActive:     input.IsActive,
		Limit:        input.PageSize,
		Offset:       input.GetOffset(),
	}
	sizeUnits, totalCount, err := s.repository.FindPaginatedSizeUnits(ctx, filter)
	if err != nil {
		return nil, err
	}
	input.SetTotalPagesAndTotalElement(totalCount)
	response := model.ResponseGetSizeUnits{
		PaginationData: &input.PaginationData,
		Data:           make([]model.SimpleSizeUnit, 0),
	}
	for _, sizeUnit := range sizeUnits {
		response.Data = append(response.Data, model.SimpleSizeUnit{
			SizeUnitID:   sizeUnit.SizeUnitID,
			SizeUnitName: sizeUnit.SizeUnitName,
			SizeUnitCode: sizeUnit.SizeUnitCode,
			SizeUnitType: sizeUnit.SizeUnitType,
			IsActive:     sizeUnit.IsActive,
			CreatedAt:    sizeUnit.CreatedAt,
			UpdatedAt:    sizeUnit.UpdatedAt,
		})
	}
	return &response, nil
}
