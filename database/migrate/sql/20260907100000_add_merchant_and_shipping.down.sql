DROP TABLE IF EXISTS order_merchant_shippings;
ALTER TABLE products DROP COLUMN IF EXISTS merchant_id;
DROP TABLE IF EXISTS merchants;