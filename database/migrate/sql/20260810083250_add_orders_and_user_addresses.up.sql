CREATE TABLE orders
(
    id             TEXT PRIMARY KEY,
    user_id        TEXT                                                         NOT NULL REFERENCES users (id),
    address_id     TEXT, -- Bisa diisi jika menggunakan tabel user_addresses
    total_amount   NUMERIC(15, 2)                                               NOT NULL,
    shipping_fee   NUMERIC(15, 2)                     DEFAULT 0.00              NOT NULL,
    status         TEXT CHECK (status IN ('PENDING', 'PAID', 'PROCESSING', 'SHIPPED', 'COMPLETED', 'CANCELLED',
                                          'EXPIRED')) DEFAULT 'PENDING'         NOT NULL,
    payment_method TEXT, -- e.g., 'midtrans', 'bank_transfer', 'credit_card'
    payment_token  TEXT, -- e.g., Snap Token Midtrans
    created_at     TIMESTAMP                          DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP

);

CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_orders_status ON orders (status);


CREATE TABLE user_addresses
(
    id             TEXT PRIMARY KEY,
    user_id        TEXT                                NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    recipient_name TEXT                                NOT NULL,
    phone_number   TEXT                                NOT NULL,
    full_address   TEXT                                NOT NULL,
    city           TEXT                                NOT NULL,
    postal_code    TEXT                                NOT NULL,
    is_primary     BOOLEAN   DEFAULT FALSE             NOT NULL,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at     TIMESTAMP
);

