-- migrate:up
ALTER TABLE products ADD COLUMN shopee_fee_free_delivery_fee DECIMAL(5,2) DEFAULT 4.00;

-- migrate:down
ALTER TABLE products DROP COLUMN shopee_fee_free_delivery_fee;
