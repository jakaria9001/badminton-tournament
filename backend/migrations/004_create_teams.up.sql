CREATE TABLE teams (
    id UUID PRIMARY KEY,

    event_id UUID NOT NULL REFERENCES events(id) ON DELETE CASCADE,

    player1_id UUID NOT NULL REFERENCES players(id),
    player2_id UUID NOT NULL REFERENCES players(id),

    team_name VARCHAR(150),

    seed INTEGER,

    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT different_players
        CHECK (player1_id <> player2_id)
);