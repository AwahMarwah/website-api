ALTER TABLE orders
    DROP COLUMN IF EXISTS expired_at,
    DROP COLUMN IF EXISTS payment_url;