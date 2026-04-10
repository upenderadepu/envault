-- +goose Up
CREATE TABLE audit_logs (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id    UUID NOT NULL REFERENCES projects(id),
    user_id       UUID REFERENCES users(id),
    action        TEXT NOT NULL,
    resource_path TEXT NOT NULL,
    ip_address    TEXT,
    user_agent    TEXT,
    request_id    TEXT,
    metadata      JSONB DEFAULT '{}',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_audit_logs_project_created ON audit_logs(project_id, created_at DESC);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);

-- Enforce immutability: the application user cannot update or delete audit rows
-- Note: REVOKE skipped for Supabase compatibility (default user is postgres/superuser)

-- +goose Down
DROP TABLE IF EXISTS audit_logs;
