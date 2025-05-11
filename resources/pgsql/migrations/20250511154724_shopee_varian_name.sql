-- migrate:up
ALTER TABLE products ADD COLUMN shopee_varian_name VARCHAR(100);

-- migrate:down
ALTER TABLE products DROP COLUMN shopee_varian_name; 
