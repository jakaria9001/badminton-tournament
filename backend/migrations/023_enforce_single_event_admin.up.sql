CREATE UNIQUE INDEX users_one_admin_per_event
    ON users (event_id)
    WHERE role = 'ADMIN' AND event_id IS NOT NULL;