ALTER TABLE merchants DROP COLUMN IF EXISTS commission_rate_bp;

DROP TABLE IF EXISTS payouts;
DROP TABLE IF EXISTS merchant_ledger_items;
DROP TABLE IF EXISTS merchant_ledgers;
