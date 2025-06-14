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
