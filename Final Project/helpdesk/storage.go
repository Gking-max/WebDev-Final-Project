package main

import (
	"context"
	"database/sql"
	"time"
)

type Model struct {
	DB *sql.DB
}

type Ticket struct {
	ID          int64
	Title       string
	Location    string
	Category    string
	Description string
	Status      string
	StatusClass string
	Priority    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SoftwareRequest struct {
	ID          int64
	Title       string
	Location    string
	RequestType string
	Description string
	CreatedAt   time.Time
}

func (m Model) InsertTicket(ctx context.Context, title, location, category, description string) (int64, error) {
	const query = `
		INSERT INTO tickets (title, location, category_id, description, status_id, priority_id)
		SELECT $1, $2, c.id, $4, s.id, p.id
		FROM categories AS c
		CROSS JOIN statuses AS s
		CROSS JOIN priorities AS p
		WHERE c.name = $3 AND s.name = 'New' AND p.name = 'Medium'
		RETURNING id`

	var id int64
	err := m.DB.QueryRowContext(ctx, query, title, location, category, description).Scan(&id)
	return id, err
}

func (m Model) GetTicket(ctx context.Context, id int64) (Ticket, error) {
	const query = `
		SELECT t.id, t.title, t.location, c.name, t.description,
		       s.name, p.name, t.created_at, t.updated_at
		FROM tickets AS t
		JOIN categories AS c ON c.id = t.category_id
		JOIN statuses AS s ON s.id = t.status_id
		JOIN priorities AS p ON p.id = t.priority_id
		WHERE t.id = $1`

	var ticket Ticket
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&ticket.ID, &ticket.Title, &ticket.Location, &ticket.Category,
		&ticket.Description, &ticket.Status, &ticket.Priority,
		&ticket.CreatedAt, &ticket.UpdatedAt,
	)
	ticket.StatusClass = statusClass(ticket.Status)
	return ticket, err
}

func (m Model) ListTickets(ctx context.Context) ([]Ticket, error) {
	const query = `
		SELECT t.id, t.title, t.location, c.name, t.description,
		       s.name, p.name, t.created_at, t.updated_at
		FROM tickets AS t
		JOIN categories AS c ON c.id = t.category_id
		JOIN statuses AS s ON s.id = t.status_id
		JOIN priorities AS p ON p.id = t.priority_id
		ORDER BY t.created_at DESC, t.id DESC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tickets := make([]Ticket, 0)
	for rows.Next() {
		var ticket Ticket
		if err := rows.Scan(
			&ticket.ID, &ticket.Title, &ticket.Location, &ticket.Category,
			&ticket.Description, &ticket.Status, &ticket.Priority,
			&ticket.CreatedAt, &ticket.UpdatedAt,
		); err != nil {
			return nil, err
		}
		ticket.StatusClass = statusClass(ticket.Status)
		tickets = append(tickets, ticket)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tickets, nil
}

func (m Model) UpdateTicket(ctx context.Context, id int64, title, location, category, description string) (int64, error) {
	const query = `
		UPDATE tickets
		SET title = $1,
		    location = $2,
		    category_id = (SELECT id FROM categories WHERE name = $3),
		    description = $4,
		    updated_at = now()
		WHERE id = $5`

	result, err := m.DB.ExecContext(ctx, query, title, location, category, description, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m Model) DeleteTicket(ctx context.Context, id int64) (int64, error) {
	result, err := m.DB.ExecContext(ctx, `DELETE FROM tickets WHERE id = $1`, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (m Model) InsertSoftwareRequest(ctx context.Context, title, location, requestType, description string) (int64, error) {
	const query = `
		INSERT INTO software_requests (title, location, description, software_request_type_id)
		SELECT $1, $2, $4, request_type.id
		FROM software_request_types AS request_type
		WHERE request_type.name = $3
		RETURNING id`

	var id int64
	err := m.DB.QueryRowContext(ctx, query, title, location, requestType, description).Scan(&id)
	return id, err
}

func (m Model) ListSoftwareRequests(ctx context.Context) ([]SoftwareRequest, error) {
	const query = `
		SELECT request.id, request.title, request.location,
		       request_type.name, request.description, request.created_at
		FROM software_requests AS request
		JOIN software_request_types AS request_type
		  ON request_type.id = request.software_request_type_id
		ORDER BY request.created_at DESC, request.id DESC`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	requests := make([]SoftwareRequest, 0)
	for rows.Next() {
		var request SoftwareRequest
		if err := rows.Scan(
			&request.ID, &request.Title, &request.Location,
			&request.RequestType, &request.Description, &request.CreatedAt,
		); err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return requests, nil
}

func statusClass(status string) string {
	switch status {
	case "New":
		return "new"
	case "Assigned":
		return "assigned"
	case "In progress":
		return "in-progress"
	case "Waiting for you":
		return "waiting"
	case "Resolved":
		return "resolved"
	default:
		return "unknown"
	}
}
