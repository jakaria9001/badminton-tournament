CREATE TABLE players (
    id UUID PRIMARY KEY,

    name VARCHAR(150) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    email VARCHAR(255),
    club_name VARCHAR(150),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_players_phone ON players(phone);