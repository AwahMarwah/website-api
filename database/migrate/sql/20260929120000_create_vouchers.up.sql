-- Voucher marketplace. merchant_id NULL berarti voucher berlaku untuk seluruh seller;
-- diisi berarti hanya berlaku untuk produk milik merchant tersebut (dicek saat checkout).

CREATE TABLE vouchers
(
    id            TEXT PRIMARY KEY,
    code          VARCHAR(50) UNIQUE NOT NULL,
    description   TEXT,
    type          VARCHAR(10)        NOT NULL CHECK (type IN ('PERCENT', 'FIXED')),
    -- PERCENT menyimpan persentase (10 = 10%), FIXED menyimpan nominal rupiah
    value         NUMERIC(15, 2)     NOT NULL CHECK (value > 0),
    max_discount  NUMERIC(15, 2),
    min_spend     NUMERIC(15, 2)     NOT NULL DEFAULT 0,
    quota         INT,
    used_count    INT                NOT NULL DEFAULT 0,
    per_user_limit INT               NOT NULL DEFAULT 1,
    starts_at     TIMESTAMP,
    ends_at       TIMESTAMP,
    merchant_id   TEXT               REFERENCES merchants (id),
    is_active     BOOLEAN            NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP,

    CHECK (ends_at IS NULL OR starts_at IS NULL OR ends_at > starts_at)
);

CREATE TABLE product_vouchers
(
    voucher_id TEXT NOT NULL REFERENCES vouchers (id) ON DELETE CASCADE,
    product_id TEXT NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    PRIMARY KEY (voucher_id, product_id)
);

-- Nominal diskon disimpan sebagai nilai jadi, bukan dihitung ulang dari voucher.
-- Voucher yang nanti diedit atau dinonaktifkan tidak boleh mengubah order lama.
CREATE TABLE voucher_redemptions
(
    id               TEXT PRIMARY KEY,
    voucher_id       TEXT          NOT NULL REFERENCES vouchers (id),
    order_id         TEXT          NOT NULL REFERENCES orders (id),
    user_id          TEXT          NOT NULL REFERENCES users (id),
    discount_amount  NUMERIC(15, 2) NOT NULL CHECK (discount_amount >= 0),
    created_at       TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_vouchers_code_active ON vouchers (code, is_active);
CREATE INDEX idx_vouchers_merchant_id ON vouchers (merchant_id);
CREATE INDEX idx_voucher_redemptions_voucher_id ON voucher_redemptions (voucher_id);
CREATE INDEX idx_voucher_redemptions_user_id ON voucher_redemptions (user_id);
CREATE INDEX idx_voucher_redemptions_order_id ON voucher_redemptions (order_id);
