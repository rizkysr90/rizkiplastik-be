package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqUpdateVarianType struct {
	*model.RequestUpdateVarianType
}

func (req *reqUpdateVarianType) sanitize() {
	req.VarianTypeID = strings.TrimSpace(req.VarianTypeID)
	req.VarianTypeName = strings.TrimSpace(strings.ToUpper(req.VarianTypeName))
	if req.VarianTypeDescription != nil {
		*req.VarianTypeDescription = strings.TrimSpace(*req.VarianTypeDescription)
	}
}
func (req *reqUpdateVarianType) validate(ctx context.Context) error {
	fieldErrors := []httperror.FieldValidation{}
	if err := common.ValidateUUIDFormat(req.VarianTypeID); err != nil {
		fieldErrors = append(fieldErrors, httperror.FieldValidation{
			Field:   "variant_type_id",
			Message: err.Error(),
		})
	}
	fieldErrors = append(fieldErrors, ValidateVarianTypeName(req.VarianTypeName)...)
	fieldErrors = append(fieldErrors, ValidateVarianTypeDescription(req.VarianTypeDescription)...)
	if len(fieldErrors) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldErrors)
	}
	return nil
}
func (s *VarianTypes) UpdateVarianType(
	ctx context.Context, data *model.RequestUpdateVarianType) error {
	input := &reqUpdateVarianType{
		RequestUpdateVarianType: data,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return err
	}
	updatedData := &repository.VarianTypeData{
		ID:        input.VarianTypeID,
		Name:      input.VarianTypeName,
		IsActive:  input.IsActive,
		UpdatedBy: "SYSTEM",
	}
	if input.VarianTypeDescription != nil {
		updatedData.Description = sql.NullString{
			String: *input.VarianTypeDescription,
			Valid:  true,
		}
	}
	if err := s.repository.UpdateTransaction(ctx, updatedData); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return httperror.NewDataNotFound(ctx, httperror.WithMessage(
				"Varian type not found",
			))
		}
		if errors.Is(err, constants.ErrAlreadyExists) {
			return httperror.NewBadRequest(ctx, httperror.WithMessage(
				"Varian type already exists",
			))
		}
		return err
	}
	return nil
}
