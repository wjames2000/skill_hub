-- Migration: 006_backfill_descriptions
-- Description: Backfill zh_description and en_description for existing skills synced before i18n columns existed

UPDATE skills SET zh_description = description WHERE zh_description IS NULL OR zh_description = '';
UPDATE skills SET en_description = description WHERE en_description IS NULL OR en_description = '';
