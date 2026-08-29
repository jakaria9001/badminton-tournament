ALTER TABLE events
ADD COLUMN best_of INTEGER NOT NULL DEFAULT 3;

ALTER TABLE events
ADD COLUMN winning_points INTEGER NOT NULL DEFAULT 21;

ALTER TABLE events
ADD COLUMN maximum_points INTEGER NOT NULL DEFAULT 30;

ALTER TABLE events
ADD CONSTRAINT chk_event_best_of
CHECK (best_of IN (1, 3));

ALTER TABLE events
ADD CONSTRAINT chk_event_winning_points
CHECK (winning_points > 0);

ALTER TABLE events
ADD CONSTRAINT chk_event_maximum_points
CHECK (maximum_points >= winning_points);