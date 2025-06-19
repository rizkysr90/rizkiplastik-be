-- migrate:up
ALTER TABLE product_categories_packaging_rules 
DROP COLUMN deleted_at,
DROP COLUMN deleted_by;


ALTER TABLE product_categories_packaging_rules ADD COLUMN is_active boolean DEFAULT true;

-- migrate:down

