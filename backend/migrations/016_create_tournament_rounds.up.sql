CREATE TABLE tournament_rounds (
    id UUID PRIMARY KEY,

    event_id UUID NOT NULL
        REFERENCES events(id)
        ON DELETE CASCADE,

    round_number INTEGER NOT NULL,

    round_name VARCHAR(40) NOT NULL,

    pairing_method VARCHAR(20) NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    locked_at TIMESTAMPTZ,

    completed_at TIMESTAMPTZ,

    UNIQUE(event_id, round_number),

    CHECK (
        pairing_method IN (
            'MANUAL',
            'RANDOM'
        )
    ),

    CHECK (
        status IN (
            'DRAFT',
            'OPEN',
            'LOCKED',
            'IN_PROGRESS',
            'COMPLETED'
        )
    )
);

CREATE INDEX idx_tournament_rounds_event
ON tournament_rounds(event_id);