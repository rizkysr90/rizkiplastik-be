package service

import (
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

func ValidateSizeUnitName(name string) []httperror.FieldValidation {
	if name == "" {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitName, "size unit name is required"),
		}
	}
	if err := common.ValidateMaxLengthStr(name, 20); err != nil {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitName, err.Error()),
		}
	}
	if err := common.ValidateMinLengthStr(name, 2); err != nil {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitName, err.Error()),
		}
	}
	return nil
}

func ValidateSizeUnitCode(code string) []httperror.FieldValidation {
	if code == "" {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitCode, "size unit code is required"),
		}
	}
	if err := common.ValidateMaxLengthStr(code, 3); err != nil {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitCode, err.Error()),
		}
	}
	if err := common.ValidateMinLengthStr(code, 1); err != nil {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitCode, err.Error()),
		}
	}
	if err := common.ValidateOnlyAllowedUppercaseLetter(code,
		fieldSizeUnitCode); err != nil {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitCode, err.Error()),
		}
	}
	return nil
}
func ValidateSizeUnitType(sizeUnitType string) []httperror.FieldValidation {
	if sizeUnitType == "" {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitType, "size unit type is required"),
		}
	}
	if err := common.ValidateEquals(sizeUnitType,
		[]string{"LENGTH", "WEIGHT", "VOLUME", "QUANTITY"}); err != nil {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitType, err.Error()),
		}
	}
	return nil
}

func ValidateSizeUnitDescription(description *string) []httperror.FieldValidation {
	if description == nil {
		return nil
	}
	if err := common.ValidateMaxLengthStr(*description, 100); err != nil {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitDescription, err.Error()),
		}
	}
	if err := common.ValidateMinLengthStr(*description, 10); err != nil {
		return []httperror.FieldValidation{
			httperror.NewFieldValidation(fieldSizeUnitDescription, err.Error()),
		}
	}
	return nil
}
