-- migrate:up
CREATE TABLE IF NOT EXISTS variant_types (
    id UUID PRIMARY KEY,
    name VARCHAR(20) NOT NULL UNIQUE,
    description VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    
    created_by VARCHAR(30) NOT NULL,
    updated_by VARCHAR(30) NOT NULL,
    
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

-- migrate:down

