CREATE TABLE registrations (
    id UUID PRIMARY KEY,

    team_id UUID NOT NULL REFERENCES teams(id) ON DELETE CASCADE,

    registered_by VARCHAR(20) NOT NULL DEFAULT 'PLAYER',

    payment_status VARCHAR(30) NOT NULL DEFAULT 'NOT_REQUIRED',

    status VARCHAR(30) NOT NULL DEFAULT 'PENDING',

    notes TEXT,

    registered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    confirmed_at TIMESTAMPTZ
);

CREATE UNIQUE INDEX idx_unique_team_registration
    ON registrations(team_id);