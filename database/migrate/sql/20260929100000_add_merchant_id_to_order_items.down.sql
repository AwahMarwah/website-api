DROP INDEX IF EXISTS idx_order_items_merchant_id;

ALTER TABLE order_items DROP COLUMN IF EXISTS merchant_id;
