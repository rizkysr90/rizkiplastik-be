package products

type CreateProductRequest struct {
	Product  Product         `json:"product"`
	Variants []VariantObject `json:"variants"`
}
