package service

import (
	"github.com/rizkysr90/rizkiplastik-be/internal/handler/variantypes/repository"
)

const (
	fVarianTypeName        = "varian_type_name"
	fVarianTypeDescription = "varian_type_description"
)

type VarianTypes struct {
	repository repository.VarianType
}

func NewVarianTypes(
	repository repository.VarianType,
) *VarianTypes {
	return &VarianTypes{
		repository: repository,
	}
}
