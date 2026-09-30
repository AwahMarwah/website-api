-- Refund. Sebelumnya notifikasi refund dari Midtrans dipetakan menjadi CANCELLED,
-- sehingga order yang uangnya sudah keluar tercampur dengan order batal sebelum bayar.
-- Status REFUNDED memisahkan keduanya.

CREATE TABLE refunds
(
    id                  TEXT PRIMARY KEY,
    order_id            TEXT          NOT NULL REFERENCES orders (id),
    merchant_id         TEXT          REFERENCES merchants (id),
    requested_by        TEXT          NOT NULL REFERENCES users (id),
    amount              NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    reason              TEXT,
    status              VARCHAR(20)    NOT NULL DEFAULT 'PENDING'
                        CHECK (status IN ('PENDING', 'APPROVED', 'PROCESSING', 'COMPLETED', 'REJECTED', 'FAILED')),
    provider            VARCHAR(30)    NOT NULL DEFAULT 'midtrans',
    provider_refund_id  TEXT,
    approved_by         TEXT,
    -- Item yang direstok, supaya bisa ditelusuri barang mana yang kembali ke gudang.
    restocked           BOOLEAN        NOT NULL DEFAULT FALSE,
    failure_reason      TEXT,
    created_at          TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP
);

CREATE INDEX idx_refunds_order_id ON refunds (order_id);
CREATE INDEX idx_refunds_status ON refunds (status, created_at DESC);
CREATE INDEX idx_refunds_merchant_id ON refunds (merchant_id);

-- Satu order tidak bisa di-refund berulang untuk alasan yang sama tanpa catatan baru.
CREATE UNIQUE INDEX idx_refunds_provider_refund_id ON refunds (provider_refund_id)
    WHERE provider_refund_id IS NOT NULL;
