DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'matches'
          AND column_name = 'loser_team_id'
    ) THEN
        ALTER TABLE matches
        ADD COLUMN loser_team_id UUID
            REFERENCES teams(id);
    END IF;
END $$;
