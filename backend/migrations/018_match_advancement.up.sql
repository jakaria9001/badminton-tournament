ALTER TABLE matches
ADD COLUMN winner_next_match_id UUID
REFERENCES matches(id);

ALTER TABLE matches
ADD COLUMN loser_next_match_id UUID
REFERENCES matches(id);

CREATE INDEX idx_matches_winner_next
ON matches(winner_next_match_id);

CREATE INDEX idx_matches_loser_next
ON matches(loser_next_match_id);