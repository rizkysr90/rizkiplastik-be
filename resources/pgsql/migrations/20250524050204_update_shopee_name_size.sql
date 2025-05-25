-- migrate:up
ALTER TABLE products 
ALTER COLUMN shopee_name TYPE VARCHAR(250);

-- migrate:down

