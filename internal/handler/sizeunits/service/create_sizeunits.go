package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqCreateSizeUnit struct {
	*model.RequestCreateSizeUnit
}

func (r *reqCreateSizeUnit) sanitize() {
	r.SizeUnitName = strings.TrimSpace(strings.ToUpper(r.SizeUnitName))
	r.SizeUnitCode = strings.TrimSpace(strings.ToUpper(r.SizeUnitCode))
	r.SizeUnitType = strings.TrimSpace(strings.ToUpper(r.SizeUnitType))
	if r.SizeUnitDescription != nil {
		*r.SizeUnitDescription = strings.TrimSpace(*r.SizeUnitDescription)
	}
}
func (r *reqCreateSizeUnit) validate(ctx context.Context) error {
	fieldValidations := []httperror.FieldValidation{}
	fieldValidations = append(fieldValidations, ValidateSizeUnitName(r.SizeUnitName)...)
	fieldValidations = append(fieldValidations, ValidateSizeUnitCode(r.SizeUnitCode)...)
	fieldValidations = append(fieldValidations, ValidateSizeUnitType(r.SizeUnitType)...)
	if r.SizeUnitDescription != nil {
		fieldValidations = append(fieldValidations, ValidateSizeUnitDescription(r.SizeUnitDescription)...)
	}
	if len(fieldValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidations)
	}
	return nil
}
func (s *SizeUnits) CreateSizeUnit(ctx context.Context,
	request model.RequestCreateSizeUnit) error {
	input := reqCreateSizeUnit{
		RequestCreateSizeUnit: &request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return err
	}
	insertedData := repository.SizeUnitData{
		SizeUnitID:   uuid.NewString(),
		SizeUnitName: input.SizeUnitName,
		SizeUnitCode: input.SizeUnitCode,
		SizeUnitType: input.SizeUnitType,
		CreatedBy:    "SYSTEM",
		UpdatedBy:    "SYSTEM",
	}
	if input.SizeUnitDescription != nil {
		insertedData.SizeUnitDescription = sql.NullString{
			String: *input.SizeUnitDescription,
			Valid:  true,
		}
	}
	if err := s.repository.InsertTransaction(ctx, insertedData); err != nil {
		if errors.Is(err, constants.ErrAlreadyExists) {
			return httperror.NewBadRequest(ctx, httperror.WithMessage(err.Error()))
		}
		return err
	}
	return nil
}
