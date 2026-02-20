-- Add missing updated_at column to government_funding_periods.
-- All other period tables (pay_plan_periods, budget_item_entries) already have it.

ALTER TABLE government_funding_periods
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMPTZ;
