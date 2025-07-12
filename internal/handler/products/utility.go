package products

import (
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

func validateFieldProduct(product *Product) []httperror.FieldValidation {
	fieldValidation := []httperror.FieldValidation{}
	if err := common.ValidateStringRequired(product.ProductType, fieldValidationFieldProductType); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldProductType,
			Message: err.Error(),
		})
	}
	if err := common.ValidateEquals(
		product.ProductType,
		[]string{
			string(repository.ProductTypeRepack),
			string(repository.ProductTypeVariant),
			string(repository.ProductTypeSingle),
		},
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldProductType,
			Message: err.Error(),
		})
	}
	fieldValidation = append(fieldValidation, ValidateBaseName(
		product.BaseName,
		fieldValidationFieldBaseName)...,
	)
	fieldValidation = append(fieldValidation, ValidateCategoryID(
		product.CategoryID,
		fieldValidationFieldCategoryID)...,
	)
	return fieldValidation
}

func validateFieldVariant(variant *VariantObject) []httperror.FieldValidation {
	fieldValidation := []httperror.FieldValidation{}
	if variant.CostPrice != nil {
		if err := common.ValidateDecimalRequired(*variant.CostPrice, fieldValidationFieldCostPrice); err != nil {
			fieldValidation = append(fieldValidation, httperror.FieldValidation{
				Field:   fieldValidationFieldCostPrice,
				Message: err.Error(),
			})
		}
	}
	if err := common.ValidateDecimalRequired(variant.SellPrice,
		fieldValidationFieldSellPrice); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldSellPrice,
			Message: err.Error(),
		})
	}
	if err := common.ValidateStringRequired(
		variant.PackagingTypeID,
		fieldValidationFieldPackagingTypeID,
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldPackagingTypeID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateStringRequired(
		variant.SizeUnitID,
		fieldValidationFieldSizeUnitID,
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldSizeUnitID,
			Message: err.Error(),
		})
	}
	if variant.SizeValue <= 0 {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldSizeValue,
			Message: "size_value must be greater than 0",
		})
	}
	if variant.CostPrice != nil && variant.CostPrice.GreaterThan(variant.SellPrice) {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldCostPrice,
			Message: "cost_price must be less than sell_price",
		})
	}
	if err := common.ValidateUUIDFormat(variant.PackagingTypeID); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldPackagingTypeID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateUUIDFormat(variant.SizeUnitID); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldSizeUnitID,
			Message: err.Error(),
		})
	}

	if variant.VariantName != nil {
		if err := common.ValidateMaxLengthStr(
			*variant.VariantName,
			constants.MaxLengthProductVariantName,
		); err != nil {
			fieldValidation = append(fieldValidation, httperror.FieldValidation{
				Field:   fieldValidationFieldVariantName,
				Message: err.Error(),
			})
		}
	}
	return fieldValidation
}
