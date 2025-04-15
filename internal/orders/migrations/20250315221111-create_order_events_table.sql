-- +migrate Up
CREATE TABLE order_events
(
    id         VARCHAR(255) PRIMARY KEY,
    order_id   VARCHAR(255) NOT NULL,
    event_type VARCHAR(255) NOT NULL,
    data       JSONB        NOT NULL    DEFAULT '{}',
    time_stamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_order_events_order_id ON order_events (order_id);
-- +migrate Down
DROP TABLE order_events;
