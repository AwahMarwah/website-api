-- Audit trail perubahan status order. Tanpa tabel ini, perubahan status adalah
-- UPDATE telanjang: tidak bisa diaudit, tidak bisa ditampilkan ke buyer sebagai timeline.

CREATE TABLE order_status_histories
(
    id         BIGSERIAL PRIMARY KEY,
    order_id   TEXT        NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    from_status VARCHAR(20),
    to_status   VARCHAR(20) NOT NULL,
    actor_id   TEXT,
    actor_role VARCHAR(30),
    note       TEXT,
    created_at TIMESTAMP   NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_osh_order_id ON order_status_histories (order_id, created_at DESC);
