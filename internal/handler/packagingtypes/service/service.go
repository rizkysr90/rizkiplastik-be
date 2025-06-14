package service

import "github.com/rizkysr90/rizkiplastik-be/internal/handler/packagingtypes/repository"

const (
	fieldPackagingID          = "packaging_id"
	fieldPackagingName        = "packaging_name"
	fieldPackagingCode        = "packaging_code"
	fieldPackagingDescription = "packaging_description"
)

// PackagingTypeService represents the packaging types service
type PackagingType struct {
	packagingTypeRepo repository.PackagingType
}

// NewPackagingType creates a new packaging types service
func NewPackagingType(packagingTypeRepo repository.PackagingType) PackagingType {
	return PackagingType{packagingTypeRepo: packagingTypeRepo}
}
