-- migrate:up
ALTER TABLE products ADD COLUMN varian_gross_profit_percentage DECIMAL(5,2);

-- migrate:down
ALTER TABLE products DROP COLUMN varian_gross_profit_percentage; 