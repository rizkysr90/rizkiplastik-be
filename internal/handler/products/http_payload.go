package products

type CreateProductRequest struct {
	Product Product
	Variant []VariantObject
}
