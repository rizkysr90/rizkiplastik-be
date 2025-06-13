package category

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
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
func (req *reqCreateCategory) validate() error {
	if err := validateCategoryName(req.Name); err != nil {
		return err
	}
	if err := common.ValidateMaxLengthStr(req.Code, 3); err != nil {
		return &util.ServiceError{
			HTTPCode: 400,
			Message:  err.Error(),
		}
	}
	if req.Description != nil {
		if err := validateCategoryDescription(*req.Description); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) CreateCategory(ctx context.Context,
	data *CreateCategoryRequest) error {
	input := reqCreateCategory{
		CreateCategoryRequest: data,
	}
	input.sanitize()
	if err := input.validate(); err != nil {
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
		return util.ConvertRepositoryError(err)
	}
	return nil
}
