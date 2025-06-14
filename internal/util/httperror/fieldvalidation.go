package httperror

import (
	"context"
	"fmt"
	"net/http"
)

type FieldValidation struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type MultiFieldValidation struct {
	Code   int               `json:"code"`
	Info   string            `json:"info"`
	Fields []FieldValidation `json:"fields"`
}

func NewFieldValidation(
	fieldName string,
	message string) FieldValidation {
	return FieldValidation{
		Field:   fieldName,
		Message: message,
	}
}
func NewMultiFieldValidation(
	_ context.Context,
	fields []FieldValidation) *MultiFieldValidation {
	return &MultiFieldValidation{
		Code:   http.StatusBadRequest,
		Info:   "INVALID_FIELD_VALIDATION",
		Fields: fields,
	}
}
func (f *MultiFieldValidation) Error() string {
	return fmt.Sprintf("MultiFieldValidation: %d - %s", f.Code, f.Info)
}
