-- migrate:up
ALTER TABLE online_transactions ADD COLUMN deleted_by VARCHAR(100);


-- migrate:down

