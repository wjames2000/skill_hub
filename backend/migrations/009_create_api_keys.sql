-- Migration: 009_create_api_keys
-- Description: Create api_keys table for user API key management

CREATE TABLE IF NOT EXISTS api_keys (
    id           BIGSERIAL PRIMARY KEY,
    user_id      BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         VARCHAR(100) NOT NULL,
    key          VARCHAR(64) NOT NULL UNIQUE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    expires_at   TIMESTAMP WITH TIME ZONE,
    is_revoked   BOOLEAN NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_api_keys_user_id ON api_keys(user_id);
CREATE INDEX IF NOT EXISTS idx_api_keys_key ON api_keys(key);
CREATE INDEX IF NOT EXISTS idx_api_keys_is_revoked ON api_keys(is_revoked);

COMMENT ON TABLE api_keys IS 'API keys for programmatic access';
COMMENT ON COLUMN api_keys.name IS 'User-friendly name for the key';
COMMENT ON COLUMN api_keys.key IS 'The actual API key value (hashed or raw)';
COMMENT ON COLUMN api_keys.is_revoked IS 'Whether the key has been revoked';
