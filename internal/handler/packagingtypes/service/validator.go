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

	// validate max length of 3
	if err := common.ValidateMaxLengthStr(code, 3); err != nil {
		fieldValidations = append(fieldValidations,
			httperror.NewFieldValidation(fieldPackagingCode, err.Error()))
	}

	// validate only allowed uppercase letter (A-Z)
	for _, char := range code {
		if char < 'A' || char > 'Z' {
			fieldValidations = append(fieldValidations,
				httperror.NewFieldValidation(fieldPackagingCode,
					"code only allowed uppercase letter (A-Z)"))
		}
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
