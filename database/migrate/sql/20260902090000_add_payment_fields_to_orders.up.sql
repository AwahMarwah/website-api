ALTER TABLE orders
    ADD COLUMN expired_at TIMESTAMP NULL,
    ADD COLUMN payment_url TEXT NULL;