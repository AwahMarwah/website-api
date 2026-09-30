-- Ledger keuangan seller. Setiap perubahan saldo punya satu baris, jadi saldo
-- selalu bisa ditelusuri ke sumbernya dan bukan hasil hitung ulang yang bisa berbeda.

CREATE TABLE merchant_ledgers
(
    id             TEXT PRIMARY KEY,
    merchant_id    TEXT          NOT NULL REFERENCES merchants (id),
    order_id       TEXT          REFERENCES orders (id),
    -- ORDER_COMMISSION: komisi marketplace yang dipotong
    -- PAYOUT: pencairan dana ke seller
    -- REFUND: koreksi karena refund
    type           VARCHAR(20)    NOT NULL CHECK (type IN ('ORDER_COMMISSION', 'PAYOUT', 'REFUND')),
    -- Amount menyimpan nilai bertanda: negatif mengurangi saldo (potongan komisi, payout, refund).
    amount         NUMERIC(15, 2) NOT NULL,
    balance_after  NUMERIC(15, 2) NOT NULL,
    note           TEXT,
    created_at     TIMESTAMP      NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_merchant_ledgers_merchant_id ON merchant_ledgers (merchant_id, created_at DESC);
CREATE INDEX idx_merchant_ledgers_order_id ON merchant_ledgers (order_id);

-- Payout tidak bisa dicatat dua kali untuk order yang sama: komisi hanya boleh
-- masuk ledger sekali per item order.
CREATE TABLE merchant_ledger_items
(
    id              TEXT PRIMARY KEY,
    ledger_id       TEXT          NOT NULL REFERENCES merchant_ledgers (id) ON DELETE CASCADE,
    order_item_id   TEXT          NOT NULL UNIQUE REFERENCES order_items (id),
    subtotal        NUMERIC(15, 2) NOT NULL,
    commission_rate_bp INTEGER    NOT NULL,
    commission      NUMERIC(15, 2) NOT NULL,
    created_at      TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE payouts
(
    id           TEXT PRIMARY KEY,
    merchant_id  TEXT          NOT NULL REFERENCES merchants (id),
    amount       NUMERIC(15, 2) NOT NULL CHECK (amount > 0),
    status       VARCHAR(20)   NOT NULL DEFAULT 'PENDING'
                 CHECK (status IN ('PENDING', 'APPROVED', 'PAID', 'REJECTED')),
    bank_account TEXT,
    note         TEXT,
    requested_by TEXT,
    approved_by  TEXT,
    ledger_id    TEXT          REFERENCES merchant_ledgers (id),
    created_at   TIMESTAMP     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP
);

CREATE INDEX idx_payouts_merchant_id ON payouts (merchant_id, created_at DESC);
CREATE INDEX idx_payouts_status ON payouts (status);

-- Rate komisi dalam basis points (1000 = 10%). Default 10% sesuai Bianchi et al. (2005)
-- untuk marketplace online.
ALTER TABLE merchants ADD COLUMN commission_rate_bp INTEGER NOT NULL DEFAULT 1000;
