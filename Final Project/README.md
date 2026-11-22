# HelpDesk

HelpDesk stores support tickets and software requests in PostgreSQL. Support
tickets have create, list, detail, edit, and delete screens; software requests
retain their create and list screens.

## Domain mapping

A HelpDesk support ticket is the project's equivalent of a report. The
assessment's `/reports` route is implemented as `/tickets`; the route names and
the database table use the HelpDesk domain consistently. The report list uses
the `.lead`, `.row-text`, and `.row-trail` row regions and the reusable `.pill`
status pattern from the list-row CSS lab.

## Database setup

Set `DATABASE_URL` in your own terminal. Do not put the real password in this
file, source control, screenshots, or submitted evidence.

```sh
export DATABASE_URL='postgres://helpdesk_jd:YOUR_PASSWORD@localhost/helpdesk_jd?sslmode=disable'
migrate -path=./migrations -database="$DATABASE_URL" up
psql "$DATABASE_URL" --file=practice.sql
```

The Go PostgreSQL driver is declared in `helpdesk/go.mod`. To download the
project's Go dependencies:

```sh
(cd helpdesk && go mod download)
```

The migrations create the lookup tables before the record tables, include
foreign-key indexes, and use `IF NOT EXISTS` for setup tables and indexes.
`practice.sql` inserts repeatable fictional sample data, including eight
tickets with varied statuses and categories.

## Run the application

```sh
(cd helpdesk && go run .)
```

The server checks `DATABASE_URL` and pings PostgreSQL before listening on
`http://localhost:4000`. The browser routes are:

| Method | Path | Purpose |
| --- | --- | --- |
| GET | `/tickets` | Newest-first list of all tickets |
| GET | `/tickets/create` | Show the create form |
| POST | `/tickets/create` | Validate and insert a ticket |
| GET | `/tickets/{id}` | View ticket details |
| GET | `/tickets/{id}/edit` | Show the prefilled edit form |
| POST | `/tickets/{id}/edit` | Validate and update the ticket |
| GET | `/tickets/{id}/delete` | Confirm deletion |
| POST | `/tickets/{id}/delete` | Delete a ticket |
| GET | `/software-requests/create` | Show the software request form |
| POST | `/software-requests/create` | Validate and insert a software request |
| GET | `/software-requests` | List software requests |

Ticket form validation is performed by the server and reinforced by database
CHECK constraints. Successful writes use Post/Redirect/Get. Unknown or invalid
ticket IDs return 404 responses. Tickets are ordered by `created_at DESC, id
DESC`; the list has no fixed limit.

## Reference adaptation

The application reuses the earlier HelpDesk app shell, forms, validation rules,
design tokens, and static assets. It adapts the class list-row structure and
status-pill CSS pattern to HelpDesk ticket fields and PostgreSQL-backed values.
The CampusFix project is treated as an implementation-flow reference, not as
the submitted application.

## Project files

- `migrations/`: paired numbered PostgreSQL up/down migrations.
- `practice.sql`: repeatable sample data and inspection queries.
- `helpdesk/`: Go HTTP application, templates, and static assets.
- `evidence.md`: catalog queries and notes for recording actual database output.
