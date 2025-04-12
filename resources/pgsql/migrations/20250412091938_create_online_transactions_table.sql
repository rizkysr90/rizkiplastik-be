-- migrate:up
-- Create enum type for online channel
CREATE TYPE ONLINE_CHANNEL AS ENUM ('SHOPEE', 'LAZADA', 'TOKOPEDIA', 'TIKTOK');

-- Create online_transactions table
CREATE TABLE online_transactions (
  id UUID PRIMARY KEY,
  type ONLINE_CHANNEL NOT NULL,
  order_number TEXT NOT NULL,
  
  created_date DATE NOT NULL,
  period_month INTEGER NOT NULL,
  period_year INTEGER NOT NULL,
  
  total_base_amount DECIMAL(10,2) NOT NULL,
  total_sale_amount DECIMAL(10,2) NOT NULL,
  total_net_profit DECIMAL(10,2) NOT NULL,
  created_by VARCHAR(30) NOT NULL,
  
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ,
  deleted_at TIMESTAMPTZ
);

-- Create online_transaction_products table
CREATE TABLE online_transaction_products (
  id UUID PRIMARY KEY,
  online_transaction_id UUID NOT NULL,
  product_name VARCHAR(50) NOT NULL,
  cost_price DECIMAL(10,2) NOT NULL,
  sale_price DECIMAL(10,2) NOT NULL,
  quantity INTEGER NOT NULL,
  
  CONSTRAINT fk_online_transaction
    FOREIGN KEY (online_transaction_id)
    REFERENCES online_transactions (id)
);

-- Add indexes
CREATE INDEX idx_online_transactions_type ON online_transactions(type);
CREATE INDEX idx_online_transactions_created_date ON online_transactions(created_date);
CREATE INDEX idx_online_transactions_period ON online_transactions(period_year, period_month);
CREATE INDEX idx_online_transaction_products_transaction_id ON online_transaction_products(online_transaction_id);

-- migrate:down

