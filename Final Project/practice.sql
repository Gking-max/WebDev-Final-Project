-- Repeatable sample data for HelpDesk. Lookup rows are inserted first.
INSERT INTO statuses (name) VALUES
  ('New'),
  ('Assigned'),
  ('In progress'),
  ('Waiting for you'),
  ('Resolved')
ON CONFLICT (name) DO NOTHING;

INSERT INTO priorities (name) VALUES
  ('Low'),
  ('Medium'),
  ('High'),
  ('Urgent')
ON CONFLICT (name) DO NOTHING;

INSERT INTO categories (name) VALUES
  ('Network'),
  ('Email'),
  ('Hardware'),
  ('Other')
ON CONFLICT (name) DO NOTHING;

INSERT INTO software_request_types (name) VALUES
  ('Installation'),
  ('Upgrade'),
  ('Licence'),
  ('Other')
ON CONFLICT (name) DO NOTHING;

WITH seed_rows(title, location, description, category, status, priority) AS (
  VALUES
    ('Wi-Fi not working in Library', 'Main Library',
     'Cannot connect to campus Wi-Fi.', 'Network', 'New', 'High'),
    ('Cannot access email', 'Student Center',
     'Getting an authentication error when signing in.', 'Email', 'Assigned', 'Medium'),
    ('Laptop running very slowly', 'Science Building 203',
     'Performance issues and a long startup time.', 'Hardware', 'In progress', 'Medium'),
    ('Cannot access student portal', 'Student Center',
     'The portal does not accept my campus login.', 'Other', 'Assigned', 'Low'),
    ('Unable to connect to campus Wi-Fi from the second-floor study area',
     'Main Library, second floor',
     'The device sees the network but cannot complete the connection.',
     'Network', 'Waiting for you', 'High'),
    ('Printer queue is stuck', 'Arts Building 110',
     'Print jobs remain queued and never reach the printer.', 'Hardware', 'New', 'Low'),
    ('Course email attachments will not open', 'Online',
     'Attachments from course messages fail to download.', 'Email', 'Resolved', 'Medium'),
    ('Shared drive unavailable in computer lab', 'Computer Lab 4',
     'The lab computers cannot open the shared course drive.', 'Other', 'New', 'Urgent')
)
INSERT INTO tickets
  (title, location, description, category_id, status_id, priority_id)
SELECT seed_rows.title, seed_rows.location, seed_rows.description,
       c.id, s.id, p.id
FROM seed_rows
JOIN categories AS c ON c.name = seed_rows.category
JOIN statuses AS s ON s.name = seed_rows.status
JOIN priorities AS p ON p.name = seed_rows.priority
WHERE NOT EXISTS (
  SELECT 1
  FROM tickets AS existing
  WHERE existing.title = seed_rows.title
    AND existing.location = seed_rows.location
);

WITH seed_rows(title, location, description, request_type) AS (
  VALUES
    ('MATLAB installation', 'Science Building 204',
     'Install MATLAB for the current course.', 'Installation'),
    ('SPSS upgrade', 'Business Lab 14',
     'Upgrade SPSS to the version required for class.', 'Upgrade')
)
INSERT INTO software_requests
  (title, location, description, software_request_type_id)
SELECT seed_rows.title, seed_rows.location, seed_rows.description, request_type.id
FROM seed_rows
JOIN software_request_types AS request_type
  ON request_type.name = seed_rows.request_type
WHERE NOT EXISTS (
  SELECT 1
  FROM software_requests AS existing
  WHERE existing.title = seed_rows.title
    AND existing.location = seed_rows.location
);

-- One query per application table.
SELECT * FROM statuses ORDER BY name;
SELECT * FROM priorities ORDER BY name;
SELECT * FROM categories ORDER BY name;
SELECT * FROM software_request_types ORDER BY name;
SELECT * FROM tickets ORDER BY created_at DESC, id DESC;
SELECT * FROM software_requests ORDER BY created_at DESC, id DESC;

-- Filter and newest-first ordering.
SELECT title, location, created_at
FROM tickets
WHERE status_id = (SELECT id FROM statuses WHERE name = 'New')
ORDER BY created_at DESC, id DESC;

SELECT id, title, created_at
FROM tickets
ORDER BY created_at DESC, id DESC;

-- Show category and status lookup values with each ticket.
SELECT t.id, t.title, t.location, c.name AS category, s.name AS status
FROM tickets AS t
JOIN categories AS c ON c.id = t.category_id
JOIN statuses AS s ON s.id = t.status_id
ORDER BY t.created_at DESC, t.id DESC;

SELECT r.id, r.title, r.location, request_type.name AS request_type
FROM software_requests AS r
JOIN software_request_types AS request_type
  ON request_type.id = r.software_request_type_id
ORDER BY r.created_at DESC, r.id DESC;

-- Grouped count by lookup value.
SELECT c.name AS category, COUNT(t.id) AS ticket_count
FROM categories AS c
LEFT JOIN tickets AS t ON t.category_id = c.id
GROUP BY c.id, c.name
ORDER BY c.name;

SELECT request_type.name AS request_type, COUNT(r.id) AS request_count
FROM software_request_types AS request_type
LEFT JOIN software_requests AS r
  ON r.software_request_type_id = request_type.id
GROUP BY request_type.id, request_type.name
ORDER BY request_type.name;
