DROP INDEX IF EXISTS idx_unique_event_player_pair;
DROP INDEX IF EXISTS idx_unique_event_player1;
DROP INDEX IF EXISTS idx_unique_event_player2;

ALTER TABLE teams
DROP COLUMN player_pair_key;