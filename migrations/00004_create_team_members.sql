-- +goose Up
CREATE TABLE team_members (
    id                   UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id           UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id              UUID NOT NULL REFERENCES users(id),
    role                 TEXT NOT NULL CHECK (role IN ('admin', 'developer', 'ci')),
    vault_policy_name    TEXT,
    vault_token_accessor TEXT,
    is_active            BOOLEAN NOT NULL DEFAULT true,
    invited_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    joined_at            TIMESTAMPTZ,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_team_members_project_user ON team_members(project_id, user_id);
CREATE INDEX idx_team_members_project_id ON team_members(project_id);
CREATE INDEX idx_team_members_user_id ON team_members(user_id);

-- +goose Down
DROP TABLE IF EXISTS team_members;
