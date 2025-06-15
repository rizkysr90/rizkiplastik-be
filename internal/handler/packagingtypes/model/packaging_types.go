package model

import "github.com/rizkysr90/rizkiplastik-be/internal/util"

type RequestCreatePackagingType struct {
	PackagingName        string  `json:"packaging_name"`
	PackagingCode        string  `json:"packaging_code"`
	PackagingDescription *string `json:"packaging_description,omitempty"`
}

type RequestUpdatePackagingType struct {
	PackagingID          string  `json:"packaging_id"`
	PackagingName        string  `json:"packaging_name"`
	PackagingDescription *string `json:"packaging_description,omitempty"`
	IsActive             bool    `json:"is_active"`
}

type RequestGetPackagingTypes struct {
	util.PaginationData `json:"pagination"`
	Name                string `json:"name"`
	Code                string `json:"code"`
	IsActive            string `json:"is_active"`
}

type ResponseGetPackagingTypes struct {
	PaginationData *util.PaginationData  `json:"pagination_data"`
	Data           []SimplePackagingType `json:"data"`
}

type RequestGetPackagingType struct {
	PackagingID string `json:"packaging_id"`
}

type ResponseGetPackagingType struct {
	Data *PackagingTypeExtended `json:"data"`
}
