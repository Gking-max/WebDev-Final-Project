CREATE TABLE IF NOT EXISTS software_request_types (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
        CHECK (name IN ('Installation', 'Upgrade', 'Licence', 'Other')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
