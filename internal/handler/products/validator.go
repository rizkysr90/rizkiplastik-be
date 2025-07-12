package products

import (
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

func ValidateBaseName(baseName, fieldName string) []httperror.FieldValidation {
	fieldValidation := []httperror.FieldValidation{}
	if err := common.ValidateStringRequired(baseName, fieldName); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldName,
			Message: err.Error(),
		})
	}
	if err := common.ValidateMaxLengthStr(
		baseName,
		constants.MaxLengthProductBaseName,
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldName,
			Message: err.Error(),
		})
	}
	return fieldValidation
}

func ValidateCategoryID(categoryID, fieldName string) []httperror.FieldValidation {
	fieldValidation := []httperror.FieldValidation{}
	if err := common.ValidateUUIDFormat(categoryID); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldName,
			Message: err.Error(),
		})
	}
	return fieldValidation
}
