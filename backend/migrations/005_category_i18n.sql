-- Migration: 005_category_i18n
-- Description: Add zh_name, en_name to skill_categories table

ALTER TABLE skill_categories ADD COLUMN IF NOT EXISTS zh_name VARCHAR(100);
ALTER TABLE skill_categories ADD COLUMN IF NOT EXISTS en_name VARCHAR(100);
