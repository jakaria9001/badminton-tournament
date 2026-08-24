CREATE TABLE events (
    id UUID PRIMARY KEY,
    tournament_id UUID NOT NULL REFERENCES tournaments(id) ON DELETE CASCADE,

    name VARCHAR(100) NOT NULL,
    event_type VARCHAR(30) NOT NULL,

    entry_fee INTEGER NOT NULL DEFAULT 0,
    max_teams INTEGER,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);