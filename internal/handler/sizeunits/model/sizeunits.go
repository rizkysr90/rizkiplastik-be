package model

import "github.com/rizkysr90/rizkiplastik-be/internal/util"

type RequestCreateSizeUnit struct {
	SizeUnitName        string  `json:"size_unit_name" binding:"required"`
	SizeUnitCode        string  `json:"size_unit_code" binding:"required"`
	SizeUnitType        string  `json:"size_unit_type" binding:"required"`
	SizeUnitDescription *string `json:"size_unit_description,omitempty"`
}

type RequestUpdateSizeUnit struct {
	SizeUnitID          string
	SizeUnitName        string  `json:"size_unit_name" binding:"required"`
	SizeUnitType        string  `json:"size_unit_type" binding:"required"`
	SizeUnitDescription *string `json:"size_unit_description,omitempty"`
	IsActive            bool    `json:"is_active"`
}
type RequestGetSizeUnits struct {
	util.PaginationData
	SizeUnitName string `json:"size_unit_name"`
	SizeUnitType string `json:"size_unit_type"`
	SizeUnitCode string `json:"size_unit_code"`
	IsActive     string `json:"is_active"`
}

type ResponseGetSizeUnits struct {
	PaginationData *util.PaginationData `json:"pagination_data"`
	Data           []SimpleSizeUnit     `json:"data"`
}
