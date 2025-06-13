package category

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type reqGetByCategoryID struct {
	*GetByCategoryIDRequest
}

func (req *reqGetByCategoryID) sanitize() {
	req.CategoryID = strings.TrimSpace(req.CategoryID)
}

func (req *reqGetByCategoryID) validate(ctx context.Context) error {
	if err := common.ValidateUUIDFormat(req.CategoryID); err != nil {
		return httperror.NewBadRequest(ctx,
			httperror.WithMessage(err.Error()))
	}
	return nil
}

func (s *Service) GetByCategoryID(
	ctx context.Context, request *GetByCategoryIDRequest) (*GetByCategoryIDResponse, error) {
	input := &reqGetByCategoryID{
		GetByCategoryIDRequest: request,
	}
	input.sanitize()
	if err := input.validate(ctx); err != nil {
		return nil, err
	}
	category, err := s.categoryRepo.GetByID(ctx, input.CategoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httperror.NewDataNotFound(ctx,
				httperror.WithMessage("category not found"))
		}
		return nil, httperror.NewInternalServer(ctx,
			httperror.WithMessage("failed to get category by id : "+err.Error()))
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
