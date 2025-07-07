package products

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rizkysr90/rizkiplastik-be/internal/repository"
)

type ProductService interface {
	Create(ctx context.Context, request *CreateProductRequest) error
}

type Service struct {
	db                       *pgxpool.Pool
	productRepository        repository.ProductRepository
	categorySizeUnitRules    repository.ProductSizeUnitRules
	categoryPackagingRules   repository.CategoryPackagingRules
	productVariantRepository repository.ProductVariant
	repackRecipeRepository   repository.RepackRecipe
}

func NewService(
	db *pgxpool.Pool,
	productRepository repository.ProductRepository,
	categorySizeUnitRules repository.ProductSizeUnitRules,
	categoryPackagingRules repository.CategoryPackagingRules,
	productVariantRepository repository.ProductVariant,
	repackRecipeRepository repository.RepackRecipe,
) ProductService {
	return &Service{
		db:                       db,
		productRepository:        productRepository,
		categorySizeUnitRules:    categorySizeUnitRules,
		categoryPackagingRules:   categoryPackagingRules,
		productVariantRepository: productVariantRepository,
		repackRecipeRepository:   repackRecipeRepository,
	}
}
