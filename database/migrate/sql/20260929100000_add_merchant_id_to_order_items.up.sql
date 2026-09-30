-- Tahap A sub-order: simpan merchant pemilik item saat transaksi dibuat.
-- Tanpa kolom ini, menentukan "order ini milik seller siapa" butuh join
-- order_items -> product_variants -> products, dan jadi rusak begitu produk
-- dipindahkan ke merchant lain.

ALTER TABLE order_items ADD COLUMN merchant_id TEXT REFERENCES merchants (id);
CREATE INDEX idx_order_items_merchant_id ON order_items (merchant_id);

-- Backfill order lama dari produk yang masih ada. Order dengan produk yang sudah
-- terhapus akan tetap NULL dan hanya bisa diperbaiki lewat adjustment manual.
UPDATE order_items oi
SET merchant_id = p.merchant_id
FROM product_variants pv
         JOIN products p ON p.id = pv.product_id
WHERE pv.id = oi.product_variant_id
  AND oi.merchant_id IS NULL
  AND p.merchant_id IS NOT NULL;
