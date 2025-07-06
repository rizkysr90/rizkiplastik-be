-- migrate:up
CREATE TABLE IF NOT EXISTS product_variants (
    id UUID PRIMARY KEY,
    product_id UUID NOT NULL REFERENCES products(id),

    product_name VARCHAR(100) NOT NULL,
    variant_name VARCHAR(100) NULL,
    full_name VARCHAR(100) NOT NULL,

    packaging_type_id UUID NOT NULL REFERENCES packaging_types(id),
    size_value DECIMAL(10, 2) NOT NULL,
    size_unit_id UUID NOT NULL REFERENCES size_units(id),

    cost_price DECIMAL(10, 2) NULL,
    selling_price DECIMAL(10, 2) NOT NULL,

    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    created_by  VARCHAR(30) NOT NULL,
    updated_by  VARCHAR(30) NOT NULL,
    deleted_by  VARCHAR(30) DEFAULT NULL,

    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz DEFAULT NULL
);
-- migrate:down

