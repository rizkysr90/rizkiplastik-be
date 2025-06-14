package model

import "time"

type SimplePackagingType struct {
	PackagingID   string    `json:"packaging_id"`
	PackagingName string    `json:"packaging_name"`
	PackagingCode string    `json:"packaging_code"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
type PackagingTypeExtended struct {
	PackagingID          string    `json:"packaging_id"`
	PackagingName        string    `json:"packaging_name"`
	PackagingCode        string    `json:"packaging_code"`
	PackagingDescription string    `json:"packaging_description"`
	IsActive             bool      `json:"is_active"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
	CreatedBy            string    `json:"created_by"`
	UpdatedBy            string    `json:"updated_by"`
}
