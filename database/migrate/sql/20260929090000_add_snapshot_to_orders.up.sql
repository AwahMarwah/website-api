-- Snapshot order: order harus withstand perubahan data master di kemudian hari.
-- address_snapshot disimpan sebagai JSONB karena bentuknya denormalisasi dari user_addresses
-- dan hanya dibaca (tidak pernah di-query per-field).

ALTER TABLE orders
    ADD COLUMN address_snapshot JSONB         NOT NULL DEFAULT '{}'::jsonb,
    ADD COLUMN buyer_name       TEXT,
    ADD COLUMN buyer_email      TEXT,
    ADD COLUMN buyer_phone      TEXT,
    ADD COLUMN note             TEXT,
    ADD COLUMN discount_amount  NUMERIC(15, 2) NOT NULL DEFAULT 0,
    ADD COLUMN paid_at          TIMESTAMP,
    ADD COLUMN shipped_at       TIMESTAMP,
    ADD COLUMN completed_at     TIMESTAMP,
    ADD COLUMN cancelled_at     TIMESTAMP;

-- Order item menyimpan nama/ sku saat transaksi terjadi. product_variants bisa di-soft delete
-- atau diubah, tapi histori order harus tetap menampilkan apa yang benar-benar dibeli.
ALTER TABLE order_items
    ADD COLUMN product_name       TEXT,
    ADD COLUMN variant_name       TEXT,
    ADD COLUMN product_image_url  TEXT,
    ADD COLUMN sku                TEXT;

CREATE INDEX idx_orders_status_created_at ON orders (status, created_at DESC);
CREATE INDEX idx_orders_user_id ON orders (user_id);
