-- migrate:up
CREATE TABLE product_error_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_number VARCHAR(255) NOT NULL,
    created_date TIMESTAMP NOT NULL,
    product_name VARCHAR(255) NOT NULL,
    error_message TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255) NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_product_error_logs_order_number ON product_error_logs(order_number);
CREATE INDEX idx_product_error_logs_created_date ON product_error_logs(created_date);
CREATE INDEX idx_product_error_logs_product_name ON product_error_logs(product_name); 

-- migrate:down
