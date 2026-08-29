CREATE TABLE round_advancements (
    id UUID PRIMARY KEY,

    round_id UUID NOT NULL
        REFERENCES tournament_rounds(id)
        ON DELETE CASCADE,

    team_id UUID NOT NULL
        REFERENCES teams(id)
        ON DELETE CASCADE,

    source_match_id UUID
        REFERENCES matches(id)
        ON DELETE SET NULL,

    advancement_type VARCHAR(20) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE(round_id, team_id),

    CHECK (
        advancement_type IN (
            'WIN',
            'BYE'
        )
    )
);

CREATE INDEX idx_round_advancements_round
ON round_advancements(round_id);