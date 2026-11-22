CREATE TABLE IF NOT EXISTS software_requests (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
        CHECK (char_length(btrim(title)) BETWEEN 1 AND 80),
    location TEXT NOT NULL
        CHECK (char_length(btrim(location)) BETWEEN 1 AND 40),
    description TEXT NOT NULL
        CHECK (char_length(btrim(description)) BETWEEN 1 AND 500),
    software_request_type_id INTEGER NOT NULL
        REFERENCES software_request_types(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_software_requests_software_request_type_id
    ON software_requests(software_request_type_id);
