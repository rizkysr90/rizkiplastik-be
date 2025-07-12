package products

import (
	"context"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/rizkysr90/rizkiplastik-be/internal/common"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
	"github.com/rizkysr90/rizkiplastik-be/internal/util/httperror"
	"github.com/shopspring/decimal"
)

type requestUpdateSingleProductType struct {
	*UpdateSingleProductTypeRequest
}

func (req *requestUpdateSingleProductType) sanitize() {
	req.PackagingTypeID = strings.TrimSpace(req.PackagingTypeID)
	req.Product.BaseName = strings.TrimSpace(req.Product.BaseName)
	req.Product.ProductType = strings.TrimSpace(req.Product.ProductType)
	req.SizeUnitID = strings.TrimSpace(req.SizeUnitID)
}

func (req *requestUpdateSingleProductType) validateField() []httperror.FieldValidation {
	fieldValidation := make([]httperror.FieldValidation, 0)
	if err := common.ValidateUUIDFormat(req.ProductID); err != nil {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldProductID,
			Message: err.Error(),
		})
	}
	if req.ProductType != string(repository.ProductTypeSingle) {
		fieldValidation = append(fieldValidation, httperror.FieldValidation{
			Field:   fieldValidationFieldProductType,
			Message: "product_type must be single",
		})
	}
	convertToProduct := Product{
		BaseName:    req.Product.BaseName,
		ProductType: req.Product.ProductType,
		CategoryID:  req.Product.CategoryID,
	}
	fieldValidation = append(fieldValidation, validateFieldProduct(&convertToProduct)...)
	convertToVariant := VariantObject{
		PackagingTypeID: req.PackagingTypeID,
		SizeValue:       req.SizeValue,
		SizeUnitID:      req.SizeUnitID,
		CostPrice:       req.CostPrice,
		SellPrice:       req.SellPrice,
		VariantName:     nil, // since simple product doesn't have variant name
		RepackRecipe:    nil, // since simple product doesn't have repack recipe
	}
	fieldValidation = append(fieldValidation, validateFieldVariant(&convertToVariant)...)
	return fieldValidation
}

func (s *Service) UpdateSingleProductType(ctx context.Context, request *UpdateSingleProductTypeRequest) error {
	userID := ctx.Value("userID").(string)
	input := &requestUpdateSingleProductType{
		UpdateSingleProductTypeRequest: request,
	}
	input.sanitize()
	fieldValidation := input.validateField()
	if len(fieldValidation) > 0 {
		return httperror.NewMultiFieldValidation(ctx, fieldValidation)
	}
	tx, err := s.db.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	variantID, err := input.validateExistingVariantData(ctx, tx, s.productVariantRepository)
	if err != nil {
		// error is already handled by validateExistingVariantData
		return err
	}
	setBaseProductUpdatedData := &repository.ProductData{
		BaseName:   input.BaseName,
		CategoryID: input.CategoryID,
		UpdatedBy:  userID,
	}
	setVariantUpdatedData := &repository.ProductVariantData{
		ID:              variantID,
		PackagingTypeID: input.PackagingTypeID,
		SizeValue:       input.SizeValue,
		SizeUnitID:      input.SizeUnitID,
		SellingPrice:    input.SellPrice,
		UpdatedBy:       userID,
		ProductName:     input.Product.BaseName,
		FullName:        input.Product.BaseName,
	}
	if input.CostPrice != nil {
		setVariantUpdatedData.CostPrice = decimal.NullDecimal{
			Decimal: *input.CostPrice, Valid: true}
	}
	if err := s.productRepository.UpdateTransaction(ctx, tx, setBaseProductUpdatedData); err != nil {
		return err
	}
	if err := s.productVariantRepository.UpdateVariantForProductTypeSingleTransaction(ctx, tx, setVariantUpdatedData); err != nil {
		return err
	}
	return nil
}
func (req *requestUpdateSingleProductType) validateExistingVariantData(ctx context.Context, tx pgx.Tx, productVariantRepository repository.ProductVariant) (string, error) {
	variants, err := productVariantRepository.FindByProductID(
		ctx, tx, req.ProductID)
	if err != nil {
		return "", err
	}
	if len(variants) != 1 {
		return "", httperror.NewBadRequest(ctx, httperror.WithMessage(
			"product_id must have exactly 1 variant : "+strconv.Itoa(len(variants)),
		))
	} else {
		return variants[0].ID, nil
	}
}
