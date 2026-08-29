ALTER TABLE matches DROP CONSTRAINT IF EXISTS matches_status_check;
ALTER TABLE registrations DROP CONSTRAINT IF EXISTS registrations_status_check;
ALTER TABLE teams DROP CONSTRAINT IF EXISTS teams_status_check;
ALTER TABLE tournaments DROP CONSTRAINT IF EXISTS tournaments_status_check;