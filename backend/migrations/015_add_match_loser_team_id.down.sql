DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'matches'
          AND column_name = 'loser_team_id'
    ) THEN
        ALTER TABLE matches
        DROP COLUMN loser_team_id;
    END IF;
END $$;
