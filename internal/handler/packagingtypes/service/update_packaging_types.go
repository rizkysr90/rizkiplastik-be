package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqUpdatePackagingType struct {
	*model.RequestUpdatePackagingType
}

func (req *reqUpdatePackagingType) sanitize() {
	req.PackagingID = strings.TrimSpace(req.PackagingID)
	req.PackagingName = strings.TrimSpace(strings.ToUpper(req.PackagingName))
	if req.PackagingDescription != nil {
		*req.PackagingDescription = strings.TrimSpace(*req.PackagingDescription)
	}
}
func (req *reqUpdatePackagingType) validateFields(ctx context.Context) error {
	fieldValidations := []httperror.FieldValidation{}
	if err := common.ValidateUUIDFormat(req.PackagingID); err != nil {
		fieldValidations = append(
			fieldValidations,
			httperror.NewFieldValidation(fieldPackagingID, err.Error()))
	}
	fieldValidations = append(fieldValidations, validateCategoryTypeName(req.PackagingName)...)
	if req.PackagingDescription != nil {
		fieldValidations = append(fieldValidations,
			validateCategoryTypeDescription(req.PackagingDescription)...)
	}
	if len(fieldValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidations)
	}
	return nil
}
func (s *PackagingType) UpdatePackagingType(
	ctx context.Context,
	data *model.RequestUpdatePackagingType) error {
	input := &reqUpdatePackagingType{
		RequestUpdatePackagingType: data,
	}
	input.sanitize()
	if err := input.validateFields(ctx); err != nil {
		return err
	}
	updatedData := &repository.PackagingTypeData{
		ID:        input.PackagingID,
		Name:      input.PackagingName,
		IsActive:  input.IsActive,
		UpdatedBy: "SYSTEM",
	}
	if input.PackagingDescription != nil {
		updatedData.Description = sql.NullString{String: *input.PackagingDescription, Valid: true}
	}
	if err := s.packagingTypeRepo.UpdateTransaction(ctx, updatedData); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return httperror.NewDataNotFound(ctx,
				httperror.WithMessage("packaging type not found"))
		}
		if errors.Is(err, constants.ErrAlreadyExists) {
			return httperror.NewBadRequest(ctx,
				httperror.WithMessage("packaging type code already exists"))
		}
		return httperror.NewInternalServer(ctx,
			httperror.WithMessage("failed to update packaging type: "+err.Error()))
	}
	return nil
}
