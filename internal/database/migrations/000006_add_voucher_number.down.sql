DROP INDEX IF EXISTS idx_child_contracts_voucher;
ALTER TABLE child_contracts DROP COLUMN IF EXISTS voucher_number;
