-- migrate:up
CREATE TABLE IF NOT EXISTS product_categories_size_unit_rules (
    rule_id UUID PRIMARY KEY,
    category_id UUID NOT NULL,
    size_unit_id UUID NOT NULL,
    
    is_default BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP NOT NULL,
    created_by VARCHAR(30) NOT NULL,
    
    updated_at TIMESTAMP NOT NULL,
    updated_by VARCHAR(30) NOT NULL,
    
    deleted_at TIMESTAMPTZ,
    deleted_by VARCHAR(30),
    
    -- Foreign key constraints
    CONSTRAINT fk_category_sizeunit_rules FOREIGN KEY (category_id) REFERENCES product_categories(id),
    CONSTRAINT fk_size_unit_sizeunit_rules FOREIGN KEY (size_unit_id) REFERENCES size_units(id)
);

CREATE UNIQUE INDEX idx_unique_sizeunit_default_per_category 
ON product_categories_size_unit_rules (category_id) 
WHERE is_default = true AND deleted_at IS NULL;



-- migrate:down

