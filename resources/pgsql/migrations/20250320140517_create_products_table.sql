-- migrate:up
CREATE TABLE IF NOT EXISTS products (
  id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  cost_price NUMERIC(10,2) NOT NULL,
  gross_profit_percentage DECIMAL(5,2),
  shopee_category VARCHAR(3),
  
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE,
  deleted_at TIMESTAMP WITH TIME ZONE
);

-- migrate:down

