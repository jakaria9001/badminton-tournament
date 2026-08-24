INSERT INTO tournaments (
    id,
    name,
    slug,
    description,
    venue_name,
    venue_address,
    start_date,
    end_date,
    status
)
VALUES (
    '00000000-0000-0000-0000-000000000001',
    'Badminton Open 2026',
    'badminton-open-2026',
    'Local Men''s Doubles Badminton Tournament',
    'Local Indoor Stadium',
    'Badarpur, Assam',
    '2026-09-12',
    '2026-09-13',
    'REGISTRATION_OPEN'
);

INSERT INTO events (
    id,
    tournament_id,
    name,
    event_type,
    entry_fee,
    max_teams
)
VALUES (
    '00000000-0000-0000-0000-000000000002',
    '00000000-0000-0000-0000-000000000001',
    'Men''s Doubles',
    'DOUBLES',
    500,
    32
);