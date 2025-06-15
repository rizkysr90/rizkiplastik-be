package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqCreateVarianType struct {
	*model.RequestCreateVarianType
}

func (req *reqCreateVarianType) sanitize() {
	req.VarianTypeName = strings.TrimSpace(strings.ToUpper(req.VarianTypeName))
	if req.VarianTypeDescription != nil {
		*req.VarianTypeDescription = strings.TrimSpace(*req.VarianTypeDescription)
	}
}
func (req *reqCreateVarianType) validate(ctx context.Context) error {
	fieldValidations := []httperror.FieldValidation{}
	fieldValidations = append(fieldValidations,
		ValidateVarianTypeName(req.VarianTypeName)...)
	fieldValidations = append(fieldValidations,
		ValidateVarianTypeDescription(req.VarianTypeDescription)...)
	if len(fieldValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidations)
	}
	return nil
}
func (s *VarianTypes) CreateVarianType(
	ctx context.Context, data *model.RequestCreateVarianType) error {
	input := &reqCreateVarianType{
		RequestCreateVarianType: data,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return err
	}
	insertedData := &repository.VarianTypeData{
		ID:        uuid.NewString(),
		Name:      input.VarianTypeName,
		CreatedBy: "SYSTEM",
		UpdatedBy: "SYSTEM",
	}
	if input.VarianTypeDescription != nil {
		insertedData.Description = sql.NullString{
			String: *input.VarianTypeDescription,
			Valid:  true,
		}
	}
	if err := s.repository.InsertTransaction(ctx, insertedData); err != nil {
		if errors.Is(err, constants.ErrAlreadyExists) {
			return httperror.NewBadRequest(ctx, httperror.WithMessage(
				"Varian type already exists",
			))
		}
		return err
	}
	return nil
}
