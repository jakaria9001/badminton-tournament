CREATE TABLE match_games (
    id UUID PRIMARY KEY,

    match_id UUID NOT NULL
        REFERENCES matches(id)
        ON DELETE CASCADE,

    game_number INTEGER NOT NULL,

    team1_score INTEGER NOT NULL DEFAULT 0,

    team2_score INTEGER NOT NULL DEFAULT 0,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(match_id, game_number),

    CHECK (team1_score >= 0),
    CHECK (team2_score >= 0)
);

CREATE INDEX idx_match_games_match
ON match_games(match_id);