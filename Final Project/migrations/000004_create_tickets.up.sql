CREATE TABLE IF NOT EXISTS tickets (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL CHECK (char_length(btrim(title)) BETWEEN 1 AND 80),
    location TEXT NOT NULL CHECK (char_length(btrim(location)) BETWEEN 1 AND 40),
    description TEXT NOT NULL CHECK (char_length(btrim(description)) BETWEEN 1 AND 500),
    status_id INTEGER NOT NULL REFERENCES statuses(id),
    priority_id INTEGER NOT NULL REFERENCES priorities(id),
    category_id INTEGER NOT NULL REFERENCES categories(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_tickets_status_id ON tickets(status_id);
CREATE INDEX IF NOT EXISTS idx_tickets_priority_id ON tickets(priority_id);
CREATE INDEX IF NOT EXISTS idx_tickets_category_id ON tickets(category_id);