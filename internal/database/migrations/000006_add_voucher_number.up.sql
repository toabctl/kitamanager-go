ALTER TABLE child_contracts ADD COLUMN voucher_number VARCHAR(50);
CREATE INDEX IF NOT EXISTS idx_child_contracts_voucher ON child_contracts(voucher_number);
