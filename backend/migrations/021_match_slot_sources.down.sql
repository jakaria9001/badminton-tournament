ALTER TABLE matches
DROP COLUMN IF EXISTS team2_source_type;

ALTER TABLE matches
DROP COLUMN IF EXISTS team2_source_match_id;

ALTER TABLE matches
DROP COLUMN IF EXISTS team1_source_type;

ALTER TABLE matches
DROP COLUMN IF EXISTS team1_source_match_id;