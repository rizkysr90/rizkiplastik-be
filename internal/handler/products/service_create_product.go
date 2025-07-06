package products

import (
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/constants"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
)

type requestCreateProduct struct {
	*CreateProductRequest
	uniqueSizeUnitArray      []string
	uniquePackagingTypeArray []string

	productCategoryCode  string
	mapSizeUnitCode      map[string]string
	mapPackagingTypeCode map[string]string
}

func (req *requestCreateProduct) sanitize() {
	sanitizeProduct(&req.Product)
	for _, variant := range req.Variant {
		sanitizeVariant(&variant)
		if variant.RepackRecipe != nil {
			sanitizeRepackRecipe(variant.RepackRecipe)
		}
	}
}
func sanitizeProduct(product *Product) {
	product.BaseName = strings.TrimSpace(product.BaseName)
	product.CategoryID = strings.TrimSpace(product.CategoryID)
	product.ProductType = strings.TrimSpace(product.ProductType)
}
func sanitizeVariant(variant *VariantObject) {
	if variant.VariantName != nil {
		*variant.VariantName = strings.TrimSpace(*variant.VariantName)
	}
	variant.PackagingTypeID = strings.TrimSpace(variant.PackagingTypeID)
	variant.SizeUnitID = strings.TrimSpace(variant.SizeUnitID)
}
func sanitizeRepackRecipe(repackRecipe *RepackRecipeObject) {
	repackRecipe.ParentVariantID = strings.TrimSpace(repackRecipe.ParentVariantID)
	repackRecipe.ChildVariantID = strings.TrimSpace(repackRecipe.ChildVariantID)
}

func (req *requestCreateProduct) validateField() []httperror.FieldValidation {
	uniqueSizeUnitID := make(map[string]bool)
	uniquePackagingTypeID := make(map[string]bool)
	fieldValidation := []httperror.FieldValidation{}
	fieldValidation = append(fieldValidation,
		validateFieldProduct(&req.Product)...)
	for _, variant := range req.Variant {
		if _, exists := uniqueSizeUnitID[variant.SizeUnitID]; !exists {
			req.uniqueSizeUnitArray = append(
				req.uniqueSizeUnitArray, variant.SizeUnitID)
		}
		if _, exists := uniquePackagingTypeID[variant.PackagingTypeID]; !exists {
			req.uniquePackagingTypeArray = append(
				req.uniquePackagingTypeArray, variant.PackagingTypeID)
		}
		uniqueSizeUnitID[variant.SizeUnitID] = true
		uniquePackagingTypeID[variant.PackagingTypeID] = true
		fieldValidation = append(fieldValidation,
			validateFieldVariant(&variant)...)
		if variant.RepackRecipe != nil {
			fieldValidation = append(fieldValidation,
				validateFieldRepackRecipe(variant.RepackRecipe)...)
		}
		if req.Product.ProductType == string(
			repository.ProductTypeRepack) && variant.RepackRecipe == nil {
			fieldValidation = append(fieldValidation, httperror.FieldValidation{
				Field:   "repack_recipe",
				Message: "repack_recipe is required for repack product",
			})
		}
	}
	return fieldValidation
}
func validateFieldProduct(product *Product) []httperror.FieldValidation {
	fieldValidation := []httperror.FieldValidation{}
	if err := common.ValidateStringRequired(product.BaseName, fieldValidationFieldBaseName); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldBaseName,
			Message: err.Error(),
		})
	}
	if err := common.ValidateStringRequired(product.CategoryID, fieldValidationFieldCategoryID); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldCategoryID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateStringRequired(product.ProductType, fieldValidationFieldProductType); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldProductType,
			Message: err.Error(),
		})
	}
	if err := common.ValidateMaxLengthStr(
		product.BaseName,
		constants.MaxLengthProductBaseName,
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldBaseName,
			Message: err.Error(),
		})
	}
	if err := common.ValidateUUIDFormat(product.CategoryID); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldCategoryID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateEquals(
		product.ProductType,
		[]string{
			string(repository.ProductTypeRepack),
			string(repository.ProductTypeVariant),
		},
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldProductType,
			Message: err.Error(),
		})
	}
	return fieldValidation
}

func validateFieldVariant(variant *VariantObject) []httperror.FieldValidation {
	fieldValidation := []httperror.FieldValidation{}
	if variant.CostPrice != nil {
		if err := common.ValidateDecimalRequired(*variant.CostPrice, fieldValidationFieldCostPrice); err != nil {
			fieldValidation = append(fieldValidation, httperror.FieldValidation{
				Field:   fieldValidationFieldCostPrice,
				Message: err.Error(),
			})
		}
	}
	if err := common.ValidateDecimalRequired(variant.SellPrice,
		fieldValidationFieldSellPrice); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldSellPrice,
			Message: err.Error(),
		})
	}
	if err := common.ValidateStringRequired(
		variant.PackagingTypeID,
		fieldValidationFieldPackagingTypeID,
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldPackagingTypeID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateStringRequired(
		variant.SizeUnitID,
		fieldValidationFieldSizeUnitID,
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldSizeUnitID,
			Message: err.Error(),
		})
	}
	if variant.SizeValue <= 0 {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldSizeValue,
			Message: "size_value must be greater than 0",
		})
	}
	if variant.CostPrice != nil && variant.CostPrice.GreaterThan(variant.SellPrice) {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldCostPrice,
			Message: "cost_price must be less than sell_price",
		})
	}
	if err := common.ValidateUUIDFormat(variant.PackagingTypeID); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldPackagingTypeID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateUUIDFormat(variant.SizeUnitID); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldSizeUnitID,
			Message: err.Error(),
		})
	}

	return fieldValidation
}
func validateFieldRepackRecipe(repackRecipe *RepackRecipeObject) []httperror.FieldValidation {
	fieldValidation := []httperror.FieldValidation{}
	if err := common.ValidateStringRequired(
		repackRecipe.ParentVariantID,
		fieldValidationFieldParentVariantID,
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldParentVariantID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateStringRequired(
		repackRecipe.ChildVariantID,
		fieldValidationFieldChildVariantID,
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldChildVariantID,
			Message: err.Error(),
		})
	}
	if err := common.ValidateDecimalRequired(
		repackRecipe.QuantityRatio,
		fieldValidationFieldQuantityRatio,
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldQuantityRatio,
			Message: err.Error(),
		})
	}
	if err := common.ValidateDecimalRequired(
		repackRecipe.RepackCostPerUnit,
		fieldValidationFieldRepackCostPerUnit,
	); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldRepackCostPerUnit,
			Message: err.Error(),
		})
	}
	if repackRecipe.RepackTimeMinutes <= 0 {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldRepackTimeMinutes,
			Message: "repack_time_minutes must be greater than 0",
		})
	}
	return fieldValidation
}
func (req *requestCreateProduct) validateCategoryRules(
	ctx context.Context,
	tx pgx.Tx,
	categoryPackagingRulesRepository repository.CategoryPackagingRules,
	categorySizeUnitRulesRepository repository.ProductSizeUnitRules,
) error {
	// Validate size unit rule
	sizeUnitRule, err := categorySizeUnitRulesRepository.FindByCategoryIDAndRuleID(
		ctx, tx, req.Product.CategoryID, req.uniqueSizeUnitArray)
	if err != nil {
		return httperror.NewInternalServer(ctx, httperror.WithMessage(
			"internal_server_error: "+err.Error(),
		))
	}
	if len(sizeUnitRule) != len(req.uniqueSizeUnitArray) {
		return httperror.NewBadRequest(ctx, httperror.WithMessage(
			"size_unit_rule_not_found",
		))
	}

	// Got product category code
	getOneSizeUnitRule := sizeUnitRule[0]
	req.productCategoryCode = getOneSizeUnitRule.ProductCategoryCode
	for _, rule := range sizeUnitRule {
		if rule.ProductCategoryCode != req.productCategoryCode {
			return httperror.NewBadRequest(ctx, httperror.WithMessage(
				"size_unit_rule_not_found",
			))
		}
		req.mapSizeUnitCode[rule.SizeUnitID] = rule.SizeUnitCode
	}
	// Validate packaging type rule
	packagingRule, err := categoryPackagingRulesRepository.FindByCategoryIDAndRuleID(
		ctx, tx, req.Product.CategoryID, req.uniquePackagingTypeArray)
	if err != nil {
		return httperror.NewInternalServer(ctx, httperror.WithMessage(
			"internal_server_error: "+err.Error(),
		))
	}
	if len(packagingRule) != len(req.uniquePackagingTypeArray) {
		return httperror.NewBadRequest(ctx, httperror.WithMessage(
			"packaging_type_rule_not_found",
		))
	}
	for _, rule := range packagingRule {
		req.mapPackagingTypeCode[rule.PackagingTypeID] = rule.PackagingTypeCode
	}
	return nil
}
func (s *Service) Create(
	ctx context.Context,
	request *CreateProductRequest) error {
	if request == nil {
		return httperror.NewBadRequest(ctx, httperror.WithMessage(
			"request is required",
		))
	}

	input := &requestCreateProduct{
		CreateProductRequest: request,
	}
	input.sanitize()
	fieldValidation := input.validateField()
	if len(fieldValidation) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidation)
	}
	// Begin database transaction
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Check if the size unit rule exists
	if err = input.validateCategoryRules(
		ctx,
		tx,
		s.categoryPackagingRules,
		s.categorySizeUnitRules,
	); err != nil {
		// error handled by function validateCategoryRules
		return err
	}
	return nil
}
