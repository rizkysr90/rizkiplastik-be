package category

import (
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

func validateCategoryName(categoryName string) []httperror.FieldValidation {
	result := []httperror.FieldValidation{}
	if err := common.ValidateMaxLengthStr(categoryName, 50); err != nil {
		result = append(result, httperror.NewFieldValidation(fieldCategoryName, err.Error()))
	}
	if err := common.ValidateMinLengthStr(categoryName, 4); err != nil {
		result = append(result, httperror.NewFieldValidation(fieldCategoryName, err.Error()))
	}
	return result
}
func validateCategoryDescription(categoryDescription string) []httperror.FieldValidation {
	result := []httperror.FieldValidation{}
	if err := common.ValidateMaxLengthStr(categoryDescription, 100); err != nil {
		result = append(result, httperror.NewFieldValidation(
			fieldCategoryDescription, err.Error()))
	}
	if err := common.ValidateMinLengthStr(categoryDescription, 10); err != nil {
		result = append(result, httperror.NewFieldValidation(
			fieldCategoryDescription, err.Error()))
	}
	return result
}
