package service

import (
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

func validateCategoryTypeName(name string) []httperror.FieldValidation {
	fieldValidations := []httperror.FieldValidation{}
	if name == "" {
		fieldValidations = append(fieldValidations,
			httperror.NewFieldValidation(fieldPackagingName, "name is required"))
	}
	if err := common.ValidateMaxLengthStr(name, 30); err != nil {
		fieldValidations = append(fieldValidations,
			httperror.NewFieldValidation(fieldPackagingName, err.Error()))
	}
	if err := common.ValidateMinLengthStr(name, 3); err != nil {
		fieldValidations = append(fieldValidations,
			httperror.NewFieldValidation(fieldPackagingName, err.Error()))
	}
	return fieldValidations
}
func validateCategoryTypeCode(code string) []httperror.FieldValidation {
	fieldValidations := []httperror.FieldValidation{}
	if code == "" {
		fieldValidations = append(fieldValidations,
			httperror.NewFieldValidation(fieldPackagingCode, "code is required"))
	}

	// Validate exact length of 1
	if err := common.ValidateLenghtEqual(code, 1); err != nil {
		fieldValidations = append(fieldValidations,
			httperror.NewFieldValidation(fieldPackagingCode, err.Error()))
	}

	// Validate that code is a single uppercase letter
	if len(code) == 1 && (code[0] < 'A' || code[0] > 'Z') {
		fieldValidations = append(fieldValidations,
			httperror.NewFieldValidation(fieldPackagingCode,
				"code must be a single uppercase letter (A-Z)"))
	}

	return fieldValidations
}
func validateCategoryTypeDescription(description *string) []httperror.FieldValidation {
	fieldValidations := []httperror.FieldValidation{}
	if description != nil {
		if err := common.ValidateMaxLengthStr(*description, 100); err != nil {
			fieldValidations = append(fieldValidations,
				httperror.NewFieldValidation(fieldPackagingDescription, err.Error()))
		}
		if err := common.ValidateMinLengthStr(*description, 10); err != nil {
			fieldValidations = append(fieldValidations,
				httperror.NewFieldValidation(fieldPackagingDescription, err.Error()))
		}
	}
	return fieldValidations
}
