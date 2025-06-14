-- migrate:up
CREATE TABLE IF NOT EXISTS packaging_types (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   name VARCHAR(30) NOT NULL,
   code CHAR(1) NOT NULL UNIQUE,
   description VARCHAR(100),
   is_active BOOLEAN DEFAULT TRUE,
   created_by VARCHAR(30) NOT NULL,
   updated_by VARCHAR(30) NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down

