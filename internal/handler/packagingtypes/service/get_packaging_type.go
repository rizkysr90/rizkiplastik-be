package service

import (
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

func (s *PackagingType) GetPackagingType(ctx context.Context,
	request *model.RequestGetPackagingType) (*model.ResponseGetPackagingType, error) {
	request.PackagingID = strings.TrimSpace(request.PackagingID)
	if err := common.ValidateUUIDFormat(request.PackagingID); err != nil {
		return nil, httperror.NewBadRequest(ctx, httperror.WithMessage(err.Error()))
	}
	packagingType, err := s.packagingTypeRepo.FindByCategoryIDExtended(ctx, request.PackagingID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, httperror.NewDataNotFound(ctx, httperror.WithMessage("packaging type not found"))
		}
		return nil, err
	}
	response := &model.ResponseGetPackagingType{
		Data: &model.PackagingTypeExtended{
			PackagingID:          packagingType.ID,
			PackagingName:        packagingType.Name,
			PackagingCode:        packagingType.Code,
			PackagingDescription: packagingType.Description.String,
			IsActive:             packagingType.IsActive,
			CreatedAt:            packagingType.CreatedAt,
			UpdatedAt:            packagingType.UpdatedAt,
			CreatedBy:            packagingType.CreatedBy,
			UpdatedBy:            packagingType.UpdatedBy,
		},
	}
	return response, nil
}
