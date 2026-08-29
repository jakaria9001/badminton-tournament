ALTER TABLE tournaments
    ADD CONSTRAINT tournaments_status_check
    CHECK (status IN ('DRAFT', 'REGISTRATION_OPEN', 'REGISTRATION_CLOSED', 'COMPLETED'));

ALTER TABLE teams
    ADD CONSTRAINT teams_status_check
    CHECK (status IN ('PENDING', 'CONFIRMED', 'REJECTED', 'WITHDRAWN'));

ALTER TABLE registrations
    ADD CONSTRAINT registrations_status_check
    CHECK (status IN ('PENDING', 'CONFIRMED', 'REJECTED', 'WITHDRAWN'));

ALTER TABLE matches
    ADD CONSTRAINT matches_status_check
    CHECK (status IN ('SCHEDULED', 'IN_PROGRESS', 'COMPLETED', 'CANCELLED'));