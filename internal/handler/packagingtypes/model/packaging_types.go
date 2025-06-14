package model

type RequestCreatePackagingType struct {
	PackagingName        string  `json:"packaging_name"`
	PackagingCode        string  `json:"packaging_code"`
	PackagingDescription *string `json:"packaging_description,omitempty"`
}
