CREATE TABLE audit_logs (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    table_name  VARCHAR(100) NOT NULL,
    record_id   UUID         NOT NULL,
    action      VARCHAR(20)  NOT NULL CHECK (action IN ('CREATE','UPDATE','DELETE')),
    user_id     UUID,
    username    VARCHAR(100) NOT NULL DEFAULT 'system',
    old_data    JSONB,
    new_data    JSONB,
    ip_address  TEXT,
    user_agent  TEXT,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_audit_table_record ON audit_logs (table_name, record_id);
CREATE INDEX idx_audit_user         ON audit_logs (user_id) WHERE user_id IS NOT NULL;
CREATE INDEX idx_audit_created      ON audit_logs (created_at DESC);
CREATE INDEX idx_audit_action       ON audit_logs (action);

-- The app role should only INSERT on this table (enforced at DB level in production).
-- REVOKE UPDATE, DELETE ON audit_logs FROM app_role;
