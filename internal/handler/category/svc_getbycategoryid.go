package category

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

type reqGetByCategoryID struct {
	*GetByCategoryIDRequest
}

func (req *reqGetByCategoryID) sanitize() {
	req.CategoryID = strings.TrimSpace(req.CategoryID)
}

func (req *reqGetByCategoryID) validate() error {
	if err := common.ValidateUUIDFormat(req.CategoryID); err != nil {
		return &util.ServiceError{
			HTTPCode: 400,
			Message:  err.Error(),
		}
	}
	return nil
}

func (s *Service) GetByCategoryID(
	ctx context.Context, request *GetByCategoryIDRequest) (*GetByCategoryIDResponse, error) {
	input := &reqGetByCategoryID{
		GetByCategoryIDRequest: request,
	}
	input.sanitize()
	if err := input.validate(); err != nil {
		return nil, err
	}
	category, err := s.categoryRepo.GetByID(ctx, input.CategoryID)
	if err != nil {
		return nil, util.ConvertRepositoryError(err)
	}
	response := &GetByCategoryIDResponse{
		CategoryDetailModel: CategoryDetailModel{
			CategoryBaseModel: CategoryBaseModel{
				CategoryID:   category.ID,
				CategoryName: category.Name,
				CategoryCode: category.Code,
				IsActive:     category.IsActive,
				CreatedAt:    category.CreatedAt,
				UpdatedAt:    category.UpdatedAt,
			},
			Description: category.Description,
		},
	}
	return response, nil
}
