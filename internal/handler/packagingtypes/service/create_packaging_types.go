package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqCreatePackagingType struct {
	*model.RequestCreatePackagingType
}

func (req *reqCreatePackagingType) sanitize() {
	req.PackagingName = strings.TrimSpace(strings.ToUpper(req.PackagingName))
	req.PackagingCode = strings.TrimSpace(strings.ToUpper(req.PackagingCode))
	if req.PackagingDescription != nil {
		description := strings.TrimSpace(*req.PackagingDescription)
		req.PackagingDescription = &description
	}
}

func (req *reqCreatePackagingType) validateFields(ctx context.Context) error {
	fieldValidations := []httperror.FieldValidation{}
	if err := validateCategoryTypeName(req.PackagingName); err != nil {
		fieldValidations = append(fieldValidations, err...)
	}
	if err := validateCategoryTypeCode(req.PackagingCode); err != nil {
		fieldValidations = append(fieldValidations, err...)
	}
	if err := validateCategoryTypeDescription(req.PackagingDescription); err != nil {
		fieldValidations = append(fieldValidations, err...)
	}
	if len(fieldValidations) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidations)
	}
	return nil
}

func (s *PackagingType) CreatePackagingType(ctx context.Context, request *model.RequestCreatePackagingType) error {
	input := &reqCreatePackagingType{
		RequestCreatePackagingType: request,
	}
	input.sanitize()
	if err := input.validateFields(ctx); err != nil {
		return err
	}

	insertData := repository.PackagingTypeData{
		ID:        uuid.NewString(),
		Name:      input.PackagingName,
		Code:      input.PackagingCode,
		CreatedBy: "SYSTEM",
		UpdatedBy: "SYSTEM",
	}

	// Handle description field
	if input.PackagingDescription != nil && *input.PackagingDescription != "" {
		insertData.Description = sql.NullString{
			String: *input.PackagingDescription,
			Valid:  true,
		}
	}

	if err := s.packagingTypeRepo.InsertTransaction(ctx, &insertData); err != nil {
		if errors.Is(err, constants.ErrAlreadyExists) {
			return httperror.NewBadRequest(
				ctx,
				httperror.WithMessage("packaging type already exists"))
		}
		return httperror.NewInternalServer(
			ctx,
			httperror.WithMessage("failed to create packaging type : "+err.Error()))
	}
	return nil
}
