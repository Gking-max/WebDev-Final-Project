# Final project evidence record

## Status before database/browser verification

The Go code and templates pass `go test ./...` and `go vet ./...`. The live
database, browser CRUD flows, responsive widths, and CSS experiments have not
been run for this record. The project database is still empty until the
migrations and seed file are applied. Record actual observations below; do not
replace them with expected results or expose the database password.

The project equivalent of a report is a HelpDesk support ticket. See
`README.md` for the domain mapping, routes, setup commands, and ordering rule.

## Database-backed list checks

| Action | Expected result | Observed result and video evidence |
| --- | --- | --- |
| Run the query below, then compare it with `/tickets`. | Every row appears, in `created_at DESC, id DESC` order, with the same ID, title, location, category, and status. | Not yet run. Record the row count, order, comparison result, and video timestamp. |
| Update a disposable test ticket's title in SQL, commit, reload `/tickets`, then restart Go. | The changed title appears without editing HTML and remains after restart. | Not yet run. Record the test ID, before/after title, restart result, and video timestamp. Restore the original title afterward. |
| Open two different ticket rows; request a missing positive ID. | Each detail page shows the ID and fields of its selected database row; the missing ID returns 404. | Not yet run. Record both IDs, titles, SQL comparison, 404 response, and video timestamp. |
| Use a disposable database for an empty-list check, or temporarily add `WHERE FALSE` to the list query and restore it. | `GET /tickets` returns 200 with “No tickets yet.” | Not yet run. Record the method used, status, empty message, and video timestamp. Do not delete working data to create this state. |
| With Go already running, stop only a local test database or introduce a controlled failing query. Restore it immediately afterward. | The page shows a generic 500 response; the server log contains diagnostic details but the browser does not show SQL errors or credentials. | Not yet run. Record the controlled failure, response, log evidence, restoration, and video timestamp. |

```sql
SELECT t.id, t.title, t.location, c.name AS category, s.name AS status,
       t.created_at
FROM tickets AS t
JOIN categories AS c ON c.id = t.category_id
JOIN statuses AS s ON s.id = t.status_id
ORDER BY t.created_at DESC, t.id DESC;
```

For the SQL title-change check, use a disposable ticket and record the original
value so it can be restored:

```sql
BEGIN;
UPDATE tickets
SET title = 'Temporary database verification title'
WHERE id = <disposable_ticket_id>;
COMMIT;
```

## CRUD verification

Use a disposable ticket and keep the eight representative seed tickets
available. Record its actual generated ID; do not assume IDs start at 1.

| Action | Expected result | Observed result and video evidence |
| --- | --- | --- |
| Submit invalid create data, then query the table. | HTTP 422; errors and submitted values are shown; no row is inserted. | Not yet run. Record response, fields, SQL row count, and video timestamp. |
| Submit valid create data. | HTTP 303 to the new ticket's details; generated ID and submitted fields match PostgreSQL. | Not yet run. Record Location header, ID, SQL result, and video timestamp. |
| Open the ticket from its list row and submit HTML-like text in a disposable field. | Details match the selected row; user text is escaped and displayed as text, not executed. | Not yet run. Record ID/title comparison and video timestamp. |
| Open edit, cancel, submit invalid values, then submit valid changes. | Cancel makes no change; invalid submission returns 422 without a database change; valid submission returns 303 and appears on details/list/SQL. The change remains after restart. | Not yet run. Record each response, before/after SQL values, restart result, and video timestamp. |
| Open and reload delete confirmation, cancel, then confirm with POST. | GET/reload and Cancel preserve the row; POST returns 303 to `/tickets`; the row is absent from the list and SQL afterward; its old detail URL returns 404. | Not yet run. Record the ID, each response, SQL result, and video timestamp. |
| Try malformed IDs and absent positive IDs on detail, edit, and delete routes. | Each record route returns 404; no other record changes and the server does not crash. | Not yet run. Record each route/status and video timestamp. |
| Observe invalid forms, successful writes, and a controlled database error. | 422 for invalid forms, 303 for successful writes, and a generic 500 for an unexpected database error. | Not yet run. Record status codes, generic error text, log evidence, and video timestamps. |

## Responsive and keyboard checks

Use the seeded long title and a multiword status. Record the actual CSS viewport
width (not only the browser window width) at each size and at 200% zoom. A
representative static row fixture was checked with the project stylesheets;
these results are preliminary and are not a substitute for checking the live
database-backed page.

| Viewport / zoom | CSS viewport width | Horizontal page overflow? | Status/title readable? | Keyboard focus and Enter navigation | Video timestamp |
| --- | --- | --- | --- | --- | --- |
| Approximately 1000px | 1000px | No (`scrollWidth` = `clientWidth` = 1000px) | Long title and status pill visible in three-region row | Not yet tested | Fixture check; no video |
| Approximately 600px | 585px (`window.innerWidth` 600px) | No (`scrollWidth` = `clientWidth` = 585px) | Pill visible; row uses narrow two-row layout | Not yet tested | Fixture check; no video |
| Approximately 360px | 345px (`window.innerWidth` 360px) | No (`scrollWidth` = `clientWidth` = 345px) | Pill visible; row uses narrow two-row layout | Tab focus visible with a 3px solid outline on row link | Fixture check; no video |
| 200% browser zoom | Not yet measured | Not yet tested | Not yet tested | Not yet tested | Not yet recorded |

## Controlled CSS experiments

Change one declaration at a time, state a prediction before the change, then
restore the final working style and verify the baseline afterward.

| Experiment | Exact change and viewport | Prediction | Observation and explanation | Baseline restored / video timestamp |
| --- | --- | --- | --- | --- |
| Grid middle track | In `.ticket-row`, compare `grid-template-columns: auto minmax(0, 1fr) auto` with `grid-template-columns: auto 1fr auto` at a narrow width, using the same long title. | Not yet written | Not yet run. Explain the observed effect of the track minimum and title wrapping; report honestly if there is no visible difference. | Not yet verified |
| Status dot `currentColor` | Temporarily change the text `color` on `.pill` or one `.pill--*` modifier. | Not yet written | Not yet run. Record whether both the label and `::before` dot change, and explain what `currentColor` refers to. | Not yet verified |

Show one relevant computed value or Grid overlay in developer tools during
these experiments. Restore the original declaration and token values.

## PostgreSQL catalog notes

Run these checks after applying migrations. Save the actual command output and
write answers based on that output.

### Tables

```sql
SELECT schemaname, tablename
FROM pg_catalog.pg_tables
WHERE schemaname = 'public'
ORDER BY tablename;
```

Expected application tables: `categories`, `priorities`,
`software_request_types`, `software_requests`, `statuses`, and `tickets`.
`schema_migrations` is maintained by `migrate`. Actual output and missing-table
answer: **not yet recorded**.

### Columns and defaults

```sql
SELECT table_name, column_name, data_type, is_nullable, column_default
FROM information_schema.columns
WHERE table_schema = 'public'
ORDER BY table_name, ordinal_position;
```

Check required record fields, lookup foreign keys, serial IDs, and `now()`
defaults for `created_at` and `updated_at`. Actual output and answer: **not yet
recorded**.

### Constraints

```sql
SELECT conrelid::regclass AS table_name, conname, contype,
       pg_catalog.pg_get_constraintdef(oid) AS definition
FROM pg_catalog.pg_constraint
WHERE connamespace = 'public'::regnamespace
ORDER BY table_name, conname;
```

Check each primary key and foreign key, the title/location/description checks,
and unique lookup names. Actual output and answer: **not yet recorded**.

### Indexes and migration state

```sql
SELECT tablename, indexname, indexdef
FROM pg_catalog.pg_indexes
WHERE schemaname = 'public'
ORDER BY tablename, indexname;

SELECT version, dirty FROM schema_migrations;
```

Verify every foreign-key column has an index. After six successful migrations,
the expected migration version is `6` and `dirty` is `false`; actual output and
answer: **not yet recorded**.

## Demonstration video

Complete video URL (Unlisted): **not yet recorded**. Confirm that the link
opens while signed out before submitting.
