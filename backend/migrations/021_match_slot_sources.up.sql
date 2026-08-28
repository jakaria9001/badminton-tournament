ALTER TABLE matches
ADD COLUMN team1_source_match_id UUID
REFERENCES matches(id);

ALTER TABLE matches
ADD COLUMN team1_source_type VARCHAR(20);

ALTER TABLE matches
ADD COLUMN team2_source_match_id UUID
REFERENCES matches(id);

ALTER TABLE matches
ADD COLUMN team2_source_type VARCHAR(20);

ALTER TABLE matches
ADD CONSTRAINT chk_team1_source_type
CHECK (
    team1_source_type IS NULL
    OR team1_source_type IN (
        'DIRECT',
        'WINNER',
        'LOSER',
        'BYE'
    )
);

ALTER TABLE matches
ADD CONSTRAINT chk_team2_source_type
CHECK (
    team2_source_type IS NULL
    OR team2_source_type IN (
        'DIRECT',
        'WINNER',
        'LOSER',
        'BYE'
    )
);