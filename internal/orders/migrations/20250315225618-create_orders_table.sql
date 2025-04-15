-- +migrate Up
CREATE TABLE orders
(
    id               VARCHAR(255) PRIMARY KEY,
    user_uuid        VARCHAR(255) NOT NULL,
    receiver_name    VARCHAR(255) NOT NULL,
    delivery_address VARCHAR(255) NOT NULL,
    created_at       TIMESTAMP WITH TIME ZONE,
    updated_at       TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_orders_user_uuid ON orders (user_uuid);
-- +migrate Down
DROP TABLE orders;
