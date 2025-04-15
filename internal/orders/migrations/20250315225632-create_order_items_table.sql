-- +migrate Up
CREATE TABLE order_items
(
    id         SERIAL PRIMARY KEY,
    order_id   VARCHAR(255) NOT NULL,
    sku_id     VARCHAR(255) NOT NULL,
    quantity   INTEGER      NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,

    ---
    CONSTRAINT fk_orders FOREIGN KEY (order_id) REFERENCES orders(id)
);

CREATE INDEX idx_order_id ON order_items(order_id);

-- +migrate Down
DROP TABLE order_items;
