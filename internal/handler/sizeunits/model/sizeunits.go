package model

type RequestCreateSizeUnit struct {
	SizeUnitName        string  `json:"size_unit_name" binding:"required"`
	SizeUnitCode        string  `json:"size_unit_code" binding:"required"`
	SizeUnitType        string  `json:"size_unit_type" binding:"required"`
	SizeUnitDescription *string `json:"size_unit_description,omitempty"`
}
