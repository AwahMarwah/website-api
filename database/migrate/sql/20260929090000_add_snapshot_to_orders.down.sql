DROP INDEX IF EXISTS idx_orders_user_id;
DROP INDEX IF EXISTS idx_orders_status_created_at;

ALTER TABLE order_items
    DROP COLUMN IF EXISTS product_image_url,
    DROP COLUMN IF EXISTS variant_name,
    DROP COLUMN IF EXISTS product_name,
    DROP COLUMN IF EXISTS sku;

ALTER TABLE orders
    DROP COLUMN IF EXISTS cancelled_at,
    DROP COLUMN IF EXISTS completed_at,
    DROP COLUMN IF EXISTS shipped_at,
    DROP COLUMN IF EXISTS paid_at,
    DROP COLUMN IF EXISTS discount_amount,
    DROP COLUMN IF EXISTS note,
    DROP COLUMN IF EXISTS buyer_phone,
    DROP COLUMN IF EXISTS buyer_email,
    DROP COLUMN IF EXISTS buyer_name,
    DROP COLUMN IF EXISTS address_snapshot;
