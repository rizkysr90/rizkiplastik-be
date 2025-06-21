-- migrate:up
ALTER TABLE product_categories_packaging_rules 
ADD CONSTRAINT fk_category_packaging_rules 
FOREIGN KEY (category_id) REFERENCES product_categories(id);

-- Add foreign key constraint for packaging_type_id
ALTER TABLE product_categories_packaging_rules 
ADD CONSTRAINT fk_packaging_type_packaging_rules 
FOREIGN KEY (packaging_type_id) REFERENCES packaging_types(id);

-- migrate:down

