ALTER TABLE teams
ADD COLUMN player_pair_key VARCHAR(100);

CREATE UNIQUE INDEX idx_unique_event_player_pair
ON teams(event_id, player_pair_key);

CREATE UNIQUE INDEX idx_unique_event_player1
ON teams(event_id, player1_id);

CREATE UNIQUE INDEX idx_unique_event_player2
ON teams(event_id, player2_id);