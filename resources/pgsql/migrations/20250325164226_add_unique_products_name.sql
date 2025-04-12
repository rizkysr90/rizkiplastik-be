-- migrate:up
ALTER TABLE products 
ADD CONSTRAINT products_name_unique UNIQUE (name);

-- migrate:down

