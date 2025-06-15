package model

import (
	"time"
)

type SimpleSizeUnit struct {
	SizeUnitID   string    `json:"size_unit_id"`
	SizeUnitName string    `json:"size_unit_name"`
	SizeUnitCode string    `json:"size_unit_code"`
	SizeUnitType string    `json:"size_unit_type"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SizeUnitExtended struct {
	SimpleSizeUnit
	SizeUnitDescription string `json:"size_unit_description"`
}
