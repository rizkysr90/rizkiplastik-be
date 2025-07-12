package products

import (
	"context"
	"database/sql"
	"log"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
	"github.com/shopspring/decimal"
)

type requestUpdateVariantProductType struct {
	*UpdateVariantProductTypeRequest
	uniqueSizeUnitArray      []string
	uniquePackagingTypeArray []string
}

func (req *requestUpdateVariantProductType) sanitizeProduct() {
	req.BaseName = strings.TrimSpace(req.BaseName)
	req.CategoryID = strings.TrimSpace(req.CategoryID)
	req.ProductID = strings.TrimSpace(req.ProductID)
}

func (req *requestUpdateVariantProductType) validateFieldProduct() []httperror.FieldValidation {
	fieldValidation := []httperror.FieldValidation{}
	if err := common.ValidateUUIDFormat(req.ProductID); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldProductID,
			Message: err.Error(),
		})
	}
	fieldValidation = append(fieldValidation, ValidateBaseName(
		req.BaseName,
		fieldValidationFieldBaseName)...,
	)
	fieldValidation = append(fieldValidation, ValidateCategoryID(
		req.CategoryID,
		fieldValidationFieldCategoryID)...,
	)
	return fieldValidation
}
func (req *requestUpdateVariantProductType) sanitizeFieldVariant(variant *VariantObject) []httperror.FieldValidation {
	variant.PackagingTypeID = strings.TrimSpace(variant.PackagingTypeID)
	variant.SizeUnitID = strings.TrimSpace(variant.SizeUnitID)
	variant.VariantID = strings.TrimSpace(variant.VariantID)
	if variant.VariantName != nil {
		*variant.VariantName = strings.TrimSpace(*variant.VariantName)
	}
	return nil
}
func (req *requestUpdateVariantProductType) validateFieldVariant(variant *VariantObject) []httperror.FieldValidation {
	fieldValidation := []httperror.FieldValidation{}
	convertToVariant := VariantObject{
		PackagingTypeID: variant.PackagingTypeID,
		SizeValue:       variant.SizeValue,
		SizeUnitID:      variant.SizeUnitID,
		CostPrice:       variant.CostPrice,
		SellPrice:       variant.SellPrice,
		VariantName:     variant.VariantName,
		RepackRecipe:    variant.RepackRecipe, // since variant product type update
	}
	fieldValidation = append(fieldValidation, validateFieldVariant(&convertToVariant)...)
	return fieldValidation
}
func (s *Service) UpdateVariantProductType(ctx context.Context,
	request *UpdateVariantProductTypeRequest) error {
	userID := ctx.Value("userID").(string)
	input := &requestUpdateVariantProductType{
		UpdateVariantProductTypeRequest: request,
	}
	input.sanitizeProduct()
	productFieldValidation := input.validateFieldProduct()
	if len(productFieldValidation) > 0 {
		return httperror.NewMultiFieldValidation(ctx, productFieldValidation)
	}
	variantFieldValidation := []httperror.FieldValidation{}
	uniqueSizeUnitID := make(map[string]bool)
	uniquePackagingTypeID := make(map[string]bool)
	for _, variant := range input.Variants {

		input.sanitizeFieldVariant(&variant)
		variantFieldValidation = append(variantFieldValidation,
			input.validateFieldVariant(&variant)...)
		if _, exists := uniqueSizeUnitID[variant.SizeUnitID]; !exists {
			input.uniqueSizeUnitArray = append(
				input.uniqueSizeUnitArray, variant.SizeUnitID)
			uniqueSizeUnitID[variant.SizeUnitID] = true

		}
		if _, exists := uniquePackagingTypeID[variant.PackagingTypeID]; !exists {
			input.uniquePackagingTypeArray = append(
				input.uniquePackagingTypeArray, variant.PackagingTypeID)
			uniquePackagingTypeID[variant.PackagingTypeID] = true
		}
	}
	if len(variantFieldValidation) > 0 {
		return httperror.NewMultiFieldValidation(ctx, variantFieldValidation)
	}
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	mapVariantIDWithData := make(map[string]*repository.ProductVariantData)
	variantProduct, err := s.productVariantRepository.FindByProductID(ctx, tx, input.ProductID)
	if err != nil {
		return err
	}
	for _, variant := range variantProduct {
		if variant.Parent.ProductType != repository.ProductTypeVariant &&
			variant.Parent.ProductType != repository.ProductTypeRepack {
			return httperror.NewBadRequest(ctx, httperror.WithMessage(
				"invalid product type variant"))
		}
		mapVariantIDWithData[variant.ID] = &variant
	}
	if len(variantProduct) != len(input.Variants) {
		return httperror.NewBadRequest(ctx, httperror.WithMessage(
			"mismatched variant product",
		))
	}
	for _, variant := range input.Variants {
		_, exist := mapVariantIDWithData[variant.VariantID]
		if !exist {
			return httperror.NewBadRequest(ctx, httperror.WithMessage(
				"product not identified : "+variant.VariantID,
			))
		}
	}
	// Validate size unit rule
	sizeUnitRule, err := s.categorySizeUnitRules.FindByCategoryIDAndSizeUnitID(
		ctx, tx, input.CategoryID, input.uniqueSizeUnitArray)
	if err != nil {
		return httperror.NewInternalServer(ctx, httperror.WithMessage(
			"internal_server_error: "+err.Error(),
		))
	}
	if len(sizeUnitRule) != len(input.uniqueSizeUnitArray) {
		log.Println("HERE len(sizeUnitRule) != len(input.uniqueSizeUnitArray)",
			len(sizeUnitRule), len(input.uniqueSizeUnitArray))
		return httperror.NewBadRequest(ctx, httperror.WithMessage(
			"size_unit_rule_not_found",
		))
	}
	// Validate packaging type rule
	packagingRule, err := s.categoryPackagingRules.FindByCategoryIDAndRuleID(
		ctx, tx, input.CategoryID, input.uniquePackagingTypeArray)
	if err != nil {
		return httperror.NewInternalServer(ctx, httperror.WithMessage(
			"internal_server_error: "+err.Error(),
		))
	}
	if len(packagingRule) != len(input.uniquePackagingTypeArray) {
		return httperror.NewBadRequest(ctx, httperror.WithMessage(
			"packaging_type_rule_not_found",
		))
	}
	setUpdatedProductData := &repository.ProductData{
		ID:         input.ProductID,
		BaseName:   strings.ToUpper(input.BaseName),
		CategoryID: input.CategoryID,
		UpdatedBy:  userID,
	}
	setUpdatedProductVariantData := make(
		[]repository.ProductVariantData, 0)
	for _, variant := range input.Variants {
		temp := repository.ProductVariantData{
			ID:          variant.VariantID,
			ProductID:   input.ProductID,
			ProductName: input.BaseName,
			// Variant name updated after temp variable
			// Full name updated after temp variable
			// Cost price updated after temp variable
			PackagingTypeID: variant.PackagingTypeID,
			SizeUnitID:      variant.SizeUnitID,
			SizeValue:       variant.SizeValue,
			SellingPrice:    variant.SellPrice,
			UpdatedBy:       userID,
		}
		if variant.VariantName != nil {
			temp.VariantName = sql.NullString{
				String: strings.ToUpper(*variant.VariantName),
				Valid:  true,
			}
			temp.FullName = setUpdatedProductData.BaseName + " " + temp.VariantName.String
		}
		if variant.CostPrice != nil {
			temp.CostPrice = decimal.NullDecimal{
				Decimal: temp.CostPrice.Decimal,
				Valid:   true,
			}
		}
		setUpdatedProductVariantData = append(
			setUpdatedProductVariantData,
			temp,
		)
	}
	if err := s.productRepository.UpdateTransaction(
		ctx, tx, setUpdatedProductData); err != nil {
		return err
	}
	for _, data := range setUpdatedProductVariantData {
		if err := s.productVariantRepository.UpdateVariantForProductTypeSingleTransaction(
			ctx, tx, &data); err != nil {
			return err
		}
	}
	if err := tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
