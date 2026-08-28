ALTER TABLE matches
ADD COLUMN round_id UUID
REFERENCES tournament_rounds(id)
ON DELETE CASCADE;

CREATE INDEX idx_matches_round
ON matches(round_id);