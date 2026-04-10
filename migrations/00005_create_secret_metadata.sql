-- +goose Up
CREATE TABLE secret_metadata (
    id               UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id       UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    environment_id   UUID NOT NULL REFERENCES environments(id) ON DELETE CASCADE,
    key_name         TEXT NOT NULL,
    vault_path       TEXT NOT NULL,
    created_by_id    UUID NOT NULL REFERENCES users(id),
    vault_version    INT NOT NULL DEFAULT 1,
    last_modified_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE UNIQUE INDEX idx_secret_metadata_unique ON secret_metadata(project_id, environment_id, key_name);
CREATE INDEX idx_secret_metadata_project_env ON secret_metadata(project_id, environment_id);

-- +goose Down
DROP TABLE IF EXISTS secret_metadata;
