package category

import (
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/util"
)

func validateCategoryName(categoryName string) error {
	if err := common.ValidateMaxLengthStr(categoryName, 50); err != nil {
		return &util.ServiceError{
			HTTPCode: 400,
			Message:  err.Error(),
		}
	}
	if err := common.ValidateMinLengthStr(categoryName, 4); err != nil {
		return &util.ServiceError{
			HTTPCode: 400,
			Message:  err.Error(),
		}
	}
	return nil
}
func validateCategoryDescription(categoryDescription string) error {
	if err := common.ValidateMaxLengthStr(categoryDescription, 100); err != nil {
		return &util.ServiceError{
			HTTPCode: 400,
			Message:  err.Error(),
		}
	}
	if err := common.ValidateMinLengthStr(categoryDescription, 10); err != nil {
		return &util.ServiceError{
			HTTPCode: 400,
			Message:  err.Error(),
		}
	}
	return nil
}
