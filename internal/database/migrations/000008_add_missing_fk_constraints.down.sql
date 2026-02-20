ALTER TABLE revoked_tokens DROP CONSTRAINT IF EXISTS fk_revoked_tokens_user;
ALTER TABLE child_attendances DROP CONSTRAINT IF EXISTS fk_child_attendances_recorded_by;
