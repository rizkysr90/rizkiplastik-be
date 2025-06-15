package service

import (
	"context"
	"strings"

	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/model"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

func (s *VarianTypes) GetVarianType(
	ctx context.Context,
	request *model.RequestGetVarianType,
) (*model.ResponseGetVarianType, error) {
	request.VarianTypeID = strings.TrimSpace(request.VarianTypeID)
	if err := common.ValidateUUIDFormat(request.VarianTypeID); err != nil {
		return nil, httperror.NewBadRequest(ctx, httperror.WithMessage(err.Error()))
	}
	variantType, err := s.repository.FindVarianTypeByIdExtended(ctx, request.VarianTypeID)
	if err != nil {
		return nil, err
	}
	response := &model.ResponseGetVarianType{
		Data: model.VarianTypeExtended{
			VarianTypeSimple: &model.VarianTypeSimple{
				VarianTypeID:   variantType.ID,
				VarianTypeName: variantType.Name,
				IsActive:       variantType.IsActive,
				CreatedAt:      variantType.CreatedAt,
				UpdatedAt:      variantType.UpdatedAt,
			},
			VarianTypeDescription: variantType.Description.String,
			CreatedBy:             variantType.CreatedBy,
			UpdatedBy:             variantType.UpdatedBy,
		},
	}
	return response, nil
}
