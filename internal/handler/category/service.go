package category

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository/pg"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type Service struct {
	categoryRepo repository.Category
}

func NewService(categoryRepo repository.Category) Service {
	return Service{
		categoryRepo: categoryRepo,
	}
}

type reqCreateCategory struct {
	*CreateCategoryRequest
}

func (req *reqCreateCategory) sanitize() {
	req.Name = strings.TrimSpace(strings.ToUpper(req.Name))
	req.Code = strings.TrimSpace(strings.ToUpper(req.Code))
	if req.Description != nil {
		*req.Description = strings.TrimSpace(strings.ToUpper(*req.Description))
	}
}
func (req *reqCreateCategory) validate(ctx context.Context) error {
	fieldValidationErrors := []httperror.FieldValidation{}
	fieldValidationErrors = append(fieldValidationErrors, validateCategoryName(req.Name)...)
	if err := common.ValidateMaxLengthStr(req.Code, 3); err != nil {
		fieldValidationErrors = append(
			fieldValidationErrors,
			httperror.NewFieldValidation(fieldCategoryCode, err.Error()))
	}
	if req.Description != nil {
		fieldValidationErrors = append(
			fieldValidationErrors,
			validateCategoryDescription(*req.Description)...)
	}
	if len(fieldValidationErrors) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidationErrors)
	}
	return nil
}

func (s *Service) CreateCategory(ctx context.Context,
	data *CreateCategoryRequest) error {
	input := reqCreateCategory{
		CreateCategoryRequest: data,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return err
	}
	// since we dont use middleware, we need to set the created by and updated by to system
	insertedData := &repository.CategoryData{
		ID:        uuid.NewString(),
		Name:      input.Name,
		Code:      input.Code,
		IsActive:  true,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if data.Description != nil {
		insertedData.Description = *data.Description
	}
	if err := s.categoryRepo.InsertTransaction(ctx, insertedData); err != nil {
		if errors.Is(err, pg.ErrAlreadyExists) {
			return httperror.NewBadRequest(
				ctx,
				httperror.WithMessage("category already exists"))
		}
		return httperror.NewInternalServer(
			ctx,
			httperror.WithMessage("failed to create category : "+err.Error()))
	}
	return nil
}
