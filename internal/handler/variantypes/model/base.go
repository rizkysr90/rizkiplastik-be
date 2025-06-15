package model

import "time"

type VarianTypeSimple struct {
	VarianTypeID   string    `json:"variant_type_id"`
	VarianTypeName string    `json:"variant_type_name"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
