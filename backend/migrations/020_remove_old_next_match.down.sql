ALTER TABLE matches
ADD COLUMN next_match_id UUID
REFERENCES matches(id);