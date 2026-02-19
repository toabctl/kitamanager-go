CREATE TABLE IF NOT EXISTS government_funding_bill_periods (
    id BIGSERIAL PRIMARY KEY,
    organization_id BIGINT NOT NULL REFERENCES organizations(id),
    from_date DATE NOT NULL,
    to_date DATE,
    file_name VARCHAR(255) NOT NULL,
    file_sha256 VARCHAR(64) NOT NULL,
    facility_name VARCHAR(255) NOT NULL,
    facility_total INT NOT NULL,
    contract_booking INT NOT NULL,
    correction_booking INT NOT NULL,
    created_by BIGINT NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_gfbp_org ON government_funding_bill_periods(organization_id);

CREATE TABLE IF NOT EXISTS government_funding_bill_children (
    id BIGSERIAL PRIMARY KEY,
    period_id BIGINT NOT NULL REFERENCES government_funding_bill_periods(id) ON DELETE CASCADE,
    voucher_number VARCHAR(20) NOT NULL,
    child_name VARCHAR(255) NOT NULL,
    birth_date VARCHAR(10) NOT NULL,
    district INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_gfbc_period ON government_funding_bill_children(period_id);

CREATE TABLE IF NOT EXISTS government_funding_bill_payments (
    id BIGSERIAL PRIMARY KEY,
    child_id BIGINT NOT NULL REFERENCES government_funding_bill_children(id) ON DELETE CASCADE,
    key VARCHAR(100) NOT NULL,
    value VARCHAR(255) NOT NULL,
    amount INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_gfbpay_child ON government_funding_bill_payments(child_id);
