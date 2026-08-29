ALTER TABLE matches
ADD COLUMN loser_team_id UUID
    REFERENCES teams(id);
