-- Add missing foreign key constraints.
-- child_attendances.recorded_by and revoked_tokens.user_id were
-- created without REFERENCES, so referential integrity was not enforced.

ALTER TABLE child_attendances
    ADD CONSTRAINT fk_child_attendances_recorded_by
    FOREIGN KEY (recorded_by) REFERENCES users(id);

ALTER TABLE revoked_tokens
    ADD CONSTRAINT fk_revoked_tokens_user
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
