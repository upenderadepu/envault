-- +goose Up
CREATE TABLE environments (
    id            UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id    UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name          TEXT NOT NULL CHECK (name IN ('development', 'staging', 'production')),
    is_production BOOLEAN NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE(project_id, name)
);

CREATE INDEX idx_environments_project_id ON environments(project_id);

-- +goose Down
DROP TABLE IF EXISTS environments;
