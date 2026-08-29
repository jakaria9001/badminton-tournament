DROP INDEX IF EXISTS idx_matches_round;

ALTER TABLE matches
DROP COLUMN round_id;