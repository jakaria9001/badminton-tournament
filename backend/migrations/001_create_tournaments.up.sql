CREATE TABLE tournaments (
    id UUID PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    slug VARCHAR(200) NOT NULL UNIQUE,
    description TEXT,

    venue_name VARCHAR(200) NOT NULL,
    venue_address TEXT,

    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    registration_deadline TIMESTAMPTZ,

    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);