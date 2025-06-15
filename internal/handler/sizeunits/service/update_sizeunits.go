package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqUpdateSizeUnit struct {
	*model.RequestUpdateSizeUnit
}

func (r *reqUpdateSizeUnit) sanitize() {
	r.SizeUnitName = strings.TrimSpace(strings.ToUpper(r.SizeUnitName))
	r.SizeUnitType = strings.TrimSpace(strings.ToUpper(r.SizeUnitType))
	if r.SizeUnitDescription != nil {
		*r.SizeUnitDescription = strings.TrimSpace(*r.SizeUnitDescription)
	}
}

func (r *reqUpdateSizeUnit) validateFields(ctx context.Context) error {
	fieldsValidations := []httperror.FieldValidation{}
	fieldsValidations = append(fieldsValidations, ValidateSizeUnitName(r.SizeUnitName)...)
	fieldsValidations = append(fieldsValidations, ValidateSizeUnitType(r.SizeUnitType)...)
	if r.SizeUnitDescription != nil {
		fieldsValidations = append(fieldsValidations,
			ValidateSizeUnitDescription(r.SizeUnitDescription)...)
	}
	if len(fieldsValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldsValidations)
	}
	return nil
}
func (s *SizeUnits) UpdateSizeUnit(ctx context.Context,
	request model.RequestUpdateSizeUnit) error {
	input := reqUpdateSizeUnit{
		RequestUpdateSizeUnit: &request,
	}
	input.sanitize()
	if err := input.validateFields(ctx); err != nil {
		return err
	}
	updatedData := repository.SizeUnitData{
		SizeUnitID:   input.SizeUnitID,
		SizeUnitName: input.SizeUnitName,
		SizeUnitType: input.SizeUnitType,
		IsActive:     input.IsActive,
		UpdatedBy:    "SYSTEM",
	}
	if input.SizeUnitDescription != nil {
		updatedData.SizeUnitDescription = sql.NullString{
			String: *input.SizeUnitDescription,
			Valid:  true,
		}
	}
	if err := s.repository.UpdateTrasaction(ctx, updatedData); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return httperror.NewDataNotFound(ctx, httperror.WithMessage("size unit not found"))
		}
		return err
	}
	return nil
}
