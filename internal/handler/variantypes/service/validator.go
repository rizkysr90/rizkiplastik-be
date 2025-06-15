package service

import (
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

func ValidateVarianTypeName(
	varianTypeName string,
) []httperror.FieldValidation {
	fieldValidations := []httperror.FieldValidation{}
	if varianTypeName == "" {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fVarianTypeName,
			Message: "Varian type name is required",
		})
	}
	if err := common.ValidateMaxLengthStr(varianTypeName, 20); err != nil {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fVarianTypeName,
			Message: err.Error(),
		})
	}
	if err := common.ValidateMinLengthStr(varianTypeName, 2); err != nil {
		fieldValidations = append(fieldValidations, httperror.FieldValidation{
			Field:   fVarianTypeName,
			Message: err.Error(),
		})
	}
	return fieldValidations
}
func ValidateVarianTypeDescription(
	varianTypeDescription *string,
) []httperror.FieldValidation {
	fieldValidations := []httperror.FieldValidation{}
	if varianTypeDescription != nil {
		if err := common.ValidateMaxLengthStr(*varianTypeDescription, 100); err != nil {
			fieldValidations = append(fieldValidations, httperror.FieldValidation{
				Field:   fVarianTypeDescription,
				Message: err.Error(),
			})
		}
		if err := common.ValidateMinLengthStr(*varianTypeDescription, 10); err != nil {
			fieldValidations = append(fieldValidations, httperror.FieldValidation{
				Field:   fVarianTypeDescription,
				Message: err.Error(),
			})
		}
	}
	return fieldValidations
}
