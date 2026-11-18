# HelpDesk — CMPS2151 Final Project

Server-rendered Go application for logging and managing IT support tickets.
The brief calls these records "reports"; this domain calls them **tickets**.
Routes, structs, templates and SQL all use `ticket` consistently.

## Domain mapping

| Brief term      | This project                                  |
|-----------------|-----------------------------------------------|
| report          | ticket                                        |
| report list     | `/tickets`                                    |
| report an issue | "New ticket" nav link                         |
| room/location   | `location_name` + `location_code`             |
| issue category  | `category_id` → `categories.name`             |

## Software

- Go 1.27
- PostgreSQL 15+
- `github.com/lib/pq`

## Setup

```bash
createdb helpdesk
export DATABASE_URL="postgres://user:password@localhost:5432/helpdesk?sslmode=disable"

psql "$DATABASE_URL" -f migrations/000001_create_statuses.up.sql
psql "$DATABASE_URL" -f migrations/000002_create_priorities.up.sql
psql "$DATABASE_URL" -f migrations/000003_create_categories.up.sql
psql "$DATABASE_URL" -f migrations/000004_create_tickets.up.sql
psql "$DATABASE_URL" -f seed.sql

go mod tidy
go run .
