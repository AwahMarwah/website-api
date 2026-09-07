CREATE TABLE merchants
(
    id             TEXT PRIMARY KEY,
    name           VARCHAR(100)        NOT NULL,
    slug           VARCHAR(100) UNIQUE NOT NULL,
    destination_id BIGINT              NOT NULL,
    city_id        VARCHAR(4),
    address        TEXT,
    is_active      BOOLEAN             DEFAULT TRUE,
    created_at     TIMESTAMP           DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP
);

ALTER TABLE products ADD COLUMN merchant_id TEXT REFERENCES merchants (id);
CREATE INDEX idx_products_merchant_id ON products (merchant_id);

CREATE TABLE order_merchant_shippings
(
    id          TEXT PRIMARY KEY,
    order_id    TEXT        NOT NULL REFERENCES orders (id),
    merchant_id TEXT        NOT NULL REFERENCES merchants (id),
    courier     VARCHAR(50) NOT NULL,
    service     VARCHAR(50),
    cost        BIGINT      NOT NULL,
    etd         VARCHAR(20),
    weight_gram INT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_oms_order_id ON order_merchant_shippings (order_id);