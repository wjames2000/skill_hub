-- Migration: 004_category_translation
-- Description: Add category_id, zh_description, en_description to skills table

ALTER TABLE skills ADD COLUMN IF NOT EXISTS category_id BIGINT REFERENCES skill_categories(id) ON DELETE SET NULL;
ALTER TABLE skills ADD COLUMN IF NOT EXISTS zh_description TEXT;
ALTER TABLE skills ADD COLUMN IF NOT EXISTS en_description TEXT;

CREATE INDEX IF NOT EXISTS idx_skills_category_id ON skills(category_id);
