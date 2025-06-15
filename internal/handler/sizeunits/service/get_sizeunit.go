package service

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

func (s *SizeUnits) GetSizeUnit(
	ctx context.Context,
	request *model.RequestGetSizeUnit) (*model.ResponseGetSizeUnit, error) {
	request.SizeUnitID = strings.TrimSpace(request.SizeUnitID)
	if err := common.ValidateUUIDFormat(request.SizeUnitID); err != nil {
		return nil, httperror.NewBadRequest(ctx, httperror.WithMessage(err.Error()))
	}
	sizeUnit, err := s.repository.FindSizeUnitByIDExtended(ctx, request.SizeUnitID)
	if err != nil {
		return nil, httperror.NewDataNotFound(ctx, httperror.WithMessage(err.Error()))
	}
	return &model.ResponseGetSizeUnit{
		Data: model.SizeUnitExtended{
			SimpleSizeUnit: model.SimpleSizeUnit{
				SizeUnitID:   sizeUnit.SizeUnitID,
				SizeUnitName: sizeUnit.SizeUnitName,
				SizeUnitCode: sizeUnit.SizeUnitCode,
				SizeUnitType: sizeUnit.SizeUnitType,
				IsActive:     sizeUnit.IsActive,
				CreatedAt:    sizeUnit.CreatedAt,
				UpdatedAt:    sizeUnit.UpdatedAt,
			},
			SizeUnitDescription: sizeUnit.SizeUnitDescription.String,
		},
	}, nil
}
