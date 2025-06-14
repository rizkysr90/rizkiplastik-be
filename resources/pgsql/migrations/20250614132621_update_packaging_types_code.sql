-- migrate:up
ALTER TABLE packaging_types 
ALTER COLUMN code TYPE varchar(3);

-- migrate:down

