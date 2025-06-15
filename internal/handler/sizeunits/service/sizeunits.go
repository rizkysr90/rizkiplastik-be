package service

import "github.com/rizkysr90/rizkiplastik-be/internal/handler/sizeunits/repository"

const (
	fieldSizeUnitName        = "size_unit_name"
	fieldSizeUnitCode        = "size_unit_code"
	fieldSizeUnitType        = "size_unit_type"
	fieldSizeUnitDescription = "size_unit_description"
)

type SizeUnits struct {
	repository repository.SizeUnits
}

func NewSizeUnits(repository repository.SizeUnits) *SizeUnits {
	return &SizeUnits{repository: repository}
}
