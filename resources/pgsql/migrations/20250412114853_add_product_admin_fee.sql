-- migrate:up
ALTER TABLE online_transaction_products
ADD COLUMN fee_amount DECIMAL(10,2) NOT NULL DEFAULT 0;

-- migrate:down

