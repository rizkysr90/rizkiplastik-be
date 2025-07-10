-- migrate:up

CREATE TABLE IF NOT EXISTS product_repack_recipes (
    id UUID PRIMARY KEY,
    parent_variant_id UUID NOT NULL REFERENCES product_variants(id),
    child_variant_id UUID NOT NULL REFERENCES product_variants(id),
    quantity_ratio DECIMAL(10, 2) NOT NULL,
    repack_cost_per_unit DECIMAL(10, 2) NOT NULL,
    repack_time_minutes INT NOT NULL,
    created_by  VARCHAR(30) NOT NULL,
    updated_by  VARCHAR(30) NOT NULL,
    deleted_by  VARCHAR(30) DEFAULT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz DEFAULT NULL
);
-- migrate:down

