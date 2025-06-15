package model

type RequestCreateVarianType struct {
	VarianTypeName        string  `json:"variant_type_name" validate:"required"`
	VarianTypeDescription *string `json:"variant_type_description,omitempty"`
}
