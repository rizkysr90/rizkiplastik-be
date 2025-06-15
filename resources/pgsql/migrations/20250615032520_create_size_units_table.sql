-- migrate:up
CREATE TABLE size_units (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   name VARCHAR(20) NOT NULL,
   code VARCHAR(3) NOT NULL UNIQUE,
   unit_type VARCHAR(10) NOT NULL,
   description VARCHAR(100),
   is_active BOOLEAN DEFAULT TRUE,
   created_by VARCHAR(30) NOT NULL,
   updated_by VARCHAR(30) NOT NULL,
   created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down

