-- +goose Up
ALTER TABLE team_members ADD COLUMN IF NOT EXISTS invite_code TEXT;
CREATE UNIQUE INDEX IF NOT EXISTS idx_team_members_invite_code ON team_members (invite_code) WHERE invite_code IS NOT NULL AND invite_code != '';

-- +goose Down
DROP INDEX IF EXISTS idx_team_members_invite_code;
ALTER TABLE team_members DROP COLUMN IF EXISTS invite_code;
