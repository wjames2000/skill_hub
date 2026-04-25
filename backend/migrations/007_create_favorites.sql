-- Migration: 007_create_favorites
-- Description: Create favorites table for user skill favorites

CREATE TABLE IF NOT EXISTS favorites (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skill_id   BIGINT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_favorites_user_id ON favorites(user_id);
CREATE INDEX IF NOT EXISTS idx_favorites_skill_id ON favorites(skill_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_favorites_user_skill ON favorites(user_id, skill_id);

COMMENT ON TABLE favorites IS 'User favorites for skills';
COMMENT ON COLUMN favorites.user_id IS 'User who favorited';
COMMENT ON COLUMN favorites.skill_id IS 'The favorited skill';
