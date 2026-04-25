-- Migration: 008_create_reviews
-- Description: Create reviews table for user skill reviews/ratings

CREATE TABLE IF NOT EXISTS reviews (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    skill_id   BIGINT NOT NULL REFERENCES skills(id) ON DELETE CASCADE,
    score      SMALLINT NOT NULL DEFAULT 0 CHECK (score >= 1 AND score <= 5),
    content    TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reviews_user_id ON reviews(user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_skill_id ON reviews(skill_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_reviews_user_skill ON reviews(user_id, skill_id);

COMMENT ON TABLE reviews IS 'User reviews and ratings for skills';
COMMENT ON COLUMN reviews.score IS 'Rating score 1-5';
COMMENT ON COLUMN reviews.content IS 'Review text content';
