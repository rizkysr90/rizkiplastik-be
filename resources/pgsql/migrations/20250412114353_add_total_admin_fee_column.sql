-- migrate:up
ALTER TABLE online_transactions
ADD COLUMN total_fee_amount DECIMAL(10,2) NOT NULL DEFAULT 0;

-- migrate:down

