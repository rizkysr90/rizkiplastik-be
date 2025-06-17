-- migrate:up
CREATE TABLE IF NOT EXISTS product_categories_packaging_rules (
    rule_id UUID PRIMARY KEY,
    category_id UUID NOT NULL,
    packaging_type_id UUID NOT NULL,
    is_default BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL,
    created_by VARCHAR(30) NOT NULL,

    updated_at TIMESTAMPTZ NOT NULL,
    updated_by VARCHAR(30) NOT NULL, 

    deleted_at TIMESTAMPTZ NULL,
    deleted_by VARCHAR(30) NULL
);

CREATE UNIQUE INDEX idx_unique_default_per_category 
ON product_categories_packaging_rules (category_id) 
WHERE is_default = TRUE AND deleted_at IS NULL;

-- migrate:down

