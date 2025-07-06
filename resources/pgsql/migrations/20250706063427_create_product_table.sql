-- migrate:up
ALTER TABLE products RENAME TO products_old;
CREATE TYPE product_type AS ENUM ('REPACK', 'VARIANT');

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY,
    base_name VARCHAR(100) NOT NULL,
    type product_type NOT NULL,
    category_id UUID NOT NULL REFERENCES product_categories(id),
    created_by  VARCHAR(30) NOT NULL,
    updated_by  VARCHAR(30) NOT NULL,
    deleted_by  VARCHAR(30) DEFAULT NULL,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz DEFAULT NULL
);
-- migrate:down

