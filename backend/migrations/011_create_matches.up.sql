CREATE TABLE matches (
    id UUID PRIMARY KEY,

    event_id UUID NOT NULL
        REFERENCES events(id)
        ON DELETE CASCADE,

    round VARCHAR(40) NOT NULL,

    match_number INTEGER NOT NULL,

    match_type VARCHAR(30) NOT NULL DEFAULT 'NORMAL',

    team1_id UUID
        REFERENCES teams(id),

    team2_id UUID
        REFERENCES teams(id),

    court_number INTEGER,

    scheduled_at TIMESTAMPTZ,

    status VARCHAR(30) NOT NULL
        DEFAULT 'SCHEDULED',

    winner_team_id UUID
        REFERENCES teams(id),

    loser_team_id UUID
        REFERENCES teams(id),

    next_match_id UUID
        REFERENCES matches(id),

    created_at TIMESTAMPTZ NOT NULL
        DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL
        DEFAULT NOW(),

    UNIQUE(event_id, round, match_number)
);

CREATE INDEX idx_matches_event
ON matches(event_id);

CREATE INDEX idx_matches_next_match
ON matches(next_match_id);

CREATE INDEX idx_matches_scheduled_at
ON matches(scheduled_at);