package main

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("record not found")

type Ticket struct {
	ID, StatusID, PriorityID, CategoryID int
	Title, Description, RequesterName, RequesterType string
	LocationName, LocationCode string
	StatusName, PriorityName, CategoryName string
	CreatedAt, UpdatedAt time.Time
}

type Option struct {
	ID   int
	Name string
}

type Model struct {
	DB *sql.DB
}

func (m *Model) ListTickets(ctx context.Context) ([]Ticket, error) {
	const q = `SELECT t.id, t.title, t.description, t.requester_name, t.requester_type,
		t.location_name, t.location_code, t.status_id, s.name, t.priority_id, p.name,
		t.category_id, c.name, t.created_at, t.updated_at
		FROM tickets t
		JOIN statuses s ON s.id = t.status_id
		JOIN priorities p ON p.id = t.priority_id
		JOIN categories c ON c.id = t.category_id
		ORDER BY t.created_at DESC, t.id DESC`
	rows, err := m.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Ticket{}
	for rows.Next() {
		var t Ticket
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.RequesterName, &t.RequesterType,
			&t.LocationName, &t.LocationCode, &t.StatusID, &t.StatusName, &t.PriorityID, &t.PriorityName,
			&t.CategoryID, &t.CategoryName, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (m *Model) GetTicket(ctx context.Context, id int) (Ticket, error) {
	const q = `SELECT t.id, t.title, t.description, t.requester_name, t.requester_type,
		t.location_name, t.location_code, t.status_id, s.name, t.priority_id, p.name,
		t.category_id, c.name, t.created_at, t.updated_at
		FROM tickets t
		JOIN statuses s ON s.id = t.status_id
		JOIN priorities p ON p.id = t.priority_id
		JOIN categories c ON c.id = t.category_id
		WHERE t.id = $1`
	var t Ticket
	err := m.DB.QueryRowContext(ctx, q, id).Scan(&t.ID, &t.Title, &t.Description, &t.RequesterName, &t.RequesterType,
		&t.LocationName, &t.LocationCode, &t.StatusID, &t.StatusName, &t.PriorityID, &t.PriorityName,
		&t.CategoryID, &t.CategoryName, &t.CreatedAt, &t.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Ticket{}, ErrNoRecord
	}
	return t, err
}

func (m *Model) InsertTicket(ctx context.Context, t *Ticket) (int, error) {
	const q = `INSERT INTO tickets (title, description, requester_name, requester_type,
		location_name, location_code, status_id, priority_id, category_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`
	var id int
	err := m.DB.QueryRowContext(ctx, q, t.Title, t.Description, t.RequesterName, t.RequesterType,
		t.LocationName, t.LocationCode, t.StatusID, t.PriorityID, t.CategoryID).Scan(&id)
	return id, err
}

func (m *Model) UpdateTicket(ctx context.Context, t *Ticket) (int64, error) {
	const q = `UPDATE tickets SET title=$1, description=$2, requester_name=$3, requester_type=$4,
		location_name=$5, location_code=$6, status_id=$7, priority_id=$8, category_id=$9, updated_at=now()
		WHERE id=$10`
	res, err := m.DB.ExecContext(ctx, q, t.Title, t.Description, t.RequesterName, t.RequesterType,
		t.LocationName, t.LocationCode, t.StatusID, t.PriorityID, t.CategoryID, t.ID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (m *Model) DeleteTicket(ctx context.Context, id int) (int64, error) {
	res, err := m.DB.ExecContext(ctx, `DELETE FROM tickets WHERE id=$1`, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (m *Model) Statuses(ctx context.Context) ([]Option, error)   { return m.options(ctx, "statuses") }
func (m *Model) Priorities(ctx context.Context) ([]Option, error) { return m.options(ctx, "priorities") }
func (m *Model) Categories(ctx context.Context) ([]Option, error) { return m.options(ctx, "categories") }

func (m *Model) options(ctx context.Context, table string) ([]Option, error) {
	rows, err := m.DB.QueryContext(ctx, "SELECT id, name FROM "+table+" ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Option{}
	for rows.Next() {
		var o Option
		if err := rows.Scan(&o.ID, &o.Name); err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, rows.Err()
}