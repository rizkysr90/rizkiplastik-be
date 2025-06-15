package model

import "github.com/rizkysr90/rizkiplastik-be/internal/util"

type RequestCreateVarianType struct {
	VarianTypeName        string  `json:"variant_type_name" validate:"required"`
	VarianTypeDescription *string `json:"variant_type_description,omitempty"`
}

type RequestUpdateVarianType struct {
	VarianTypeID          string
	VarianTypeName        string  `json:"variant_type_name" validate:"required"`
	VarianTypeDescription *string `json:"variant_type_description,omitempty"`
	IsActive              bool    `json:"is_active"`
}

type RequestVarianTypePaginated struct {
	util.PaginationData
	VarianTypeName string `json:"variant_type_name"`
	IsActive       string `json:"is_active"`
}

type ResponseVarianTypePaginated struct {
	Pagination *util.PaginationData `json:"pagination"`
	Data       []VarianTypeSimple   `json:"data"`
}

type RequestGetVarianType struct {
	VarianTypeID string `json:"variant_type_id"`
}

type ResponseGetVarianType struct {
	Data VarianTypeExtended `json:"data"`
}
