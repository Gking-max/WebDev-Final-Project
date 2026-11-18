package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type App struct {
	Model *Model
	Cache map[string]*template.Template
}

func (a *App) render(w http.ResponseWriter, status int, page string, data any) {
	ts, ok := a.Cache[page]
	if !ok {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	var buf bytes.Buffer
	if err := ts.ExecuteTemplate(&buf, "base", data); err != nil {
		log.Printf("render %q: %v", page, err)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	buf.WriteTo(w)
}

func (a *App) serverError(w http.ResponseWriter, err error) {
	log.Printf("server error: %v", err)
	http.Error(w, "Internal Server Error", 500)
}

func (a *App) notFound(w http.ResponseWriter) {
	http.Error(w, "404 page not found", http.StatusNotFound)
}

func parseID(r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		return 0, false
	}
	return id, true
}

func statusMod(s string) string {
	return strings.ToLower(strings.ReplaceAll(strings.TrimSpace(s), " ", "-"))
}

func (a *App) home(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/tickets", http.StatusSeeOther)
}

func (a *App) listTickets(w http.ResponseWriter, r *http.Request) {
	tickets, err := a.Model.ListTickets(r.Context())
	if err != nil {
		a.serverError(w, err)
		return
	}
	a.render(w, http.StatusOK, "list.html", map[string]any{"Tickets": tickets, "Active": "tickets"})
}

func (a *App) showTicket(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		a.notFound(w)
		return
	}
	t, err := a.Model.GetTicket(r.Context(), id)
	if errors.Is(err, ErrNoRecord) {
		a.notFound(w)
		return
	}
	if err != nil {
		a.serverError(w, err)
		return
	}
	a.render(w, http.StatusOK, "details.html", map[string]any{"Ticket": t, "Active": "tickets"})
}

func (a *App) createForm(w http.ResponseWriter, r *http.Request) {
	a.renderForm(w, r, http.StatusOK, Ticket{}, map[string]string{}, true)
}

func (a *App) createSubmit(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		a.serverError(w, err)
		return
	}
	t, errs := a.parseTicketForm(r)
	if len(errs) > 0 {
		a.renderForm(w, r, http.StatusUnprocessableEntity, t, errs, true)
		return
	}
	id, err := a.Model.InsertTicket(r.Context(), &t)
	if err != nil {
		a.serverError(w, err)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/tickets/%d", id), http.StatusSeeOther)
}

func (a *App) editForm(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		a.notFound(w)
		return
	}
	t, err := a.Model.GetTicket(r.Context(), id)
	if errors.Is(err, ErrNoRecord) {
		a.notFound(w)
		return
	}
	if err != nil {
		a.serverError(w, err)
		return
	}
	a.renderForm(w, r, http.StatusOK, t, map[string]string{}, false)
}

func (a *App) editSubmit(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		a.notFound(w)
		return
	}
	if _, err := a.Model.GetTicket(r.Context(), id); errors.Is(err, ErrNoRecord) {
		a.notFound(w)
		return
	} else if err != nil {
		a.serverError(w, err)
		return
	}
	if err := r.ParseForm(); err != nil {
		a.serverError(w, err)
		return
	}
	t, errs := a.parseTicketForm(r)
	t.ID = id
	if len(errs) > 0 {
		a.renderForm(w, r, http.StatusUnprocessableEntity, t, errs, false)
		return
	}
	rows, err := a.Model.UpdateTicket(r.Context(), &t)
	if err != nil {
		a.serverError(w, err)
		return
	}
	if rows == 0 {
		a.notFound(w)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/tickets/%d", id), http.StatusSeeOther)
}

func (a *App) deleteConfirm(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		a.notFound(w)
		return
	}
	t, err := a.Model.GetTicket(r.Context(), id)
	if errors.Is(err, ErrNoRecord) {
		a.notFound(w)
		return
	}
	if err != nil {
		a.serverError(w, err)
		return
	}
	a.render(w, http.StatusOK, "delete.html", map[string]any{"Ticket": t, "Active": "tickets"})
}

func (a *App) deleteSubmit(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(r)
	if !ok {
		a.notFound(w)
		return
	}
	rows, err := a.Model.DeleteTicket(r.Context(), id)
	if err != nil {
		a.serverError(w, err)
		return
	}
	if rows == 0 {
		a.notFound(w)
		return
	}
	http.Redirect(w, r, "/tickets", http.StatusSeeOther)
}

func (a *App) renderForm(w http.ResponseWriter, r *http.Request, status int, t Ticket, errs map[string]string, isCreate bool) {
	statuses, err := a.Model.Statuses(r.Context())
	if err != nil {
		a.serverError(w, err)
		return
	}
	priorities, err := a.Model.Priorities(r.Context())
	if err != nil {
		a.serverError(w, err)
		return
	}
	categories, err := a.Model.Categories(r.Context())
	if err != nil {
		a.serverError(w, err)
		return
	}
	action := "/tickets/create"
	if !isCreate {
		action = fmt.Sprintf("/tickets/%d/edit", t.ID)
	}
	a.render(w, status, "form.html", map[string]any{
		"Ticket": t, "Errors": errs, "IsCreate": isCreate, "Action": action,
		"Statuses": statuses, "Priorities": priorities, "Categories": categories, "Active": "tickets",
	})
}

func (a *App) parseTicketForm(r *http.Request) (Ticket, map[string]string) {
	errs := map[string]string{}
	t := Ticket{
		Title:         strings.TrimSpace(r.PostFormValue("title")),
		Description:   strings.TrimSpace(r.PostFormValue("description")),
		RequesterName: strings.TrimSpace(r.PostFormValue("requester_name")),
		RequesterType: r.PostFormValue("requester_type"),
		LocationName:  strings.TrimSpace(r.PostFormValue("location_name")),
		LocationCode:  strings.TrimSpace(r.PostFormValue("location_code")),
	}
	t.StatusID, _ = strconv.Atoi(r.PostFormValue("status_id"))
	t.PriorityID, _ = strconv.Atoi(r.PostFormValue("priority_id"))
	t.CategoryID, _ = strconv.Atoi(r.PostFormValue("category_id"))

	if t.Title == "" {
		errs["title"] = "Title is required."
	}
	if t.Description == "" {
		errs["description"] = "Description is required."
	}
	if t.RequesterName == "" {
		errs["requester_name"] = "Requester name is required."
	}
	if t.RequesterType != "Student" && t.RequesterType != "Faculty" && t.RequesterType != "Staff" {
		errs["requester_type"] = "Choose Student, Faculty or Staff."
	}
	if t.LocationName == "" {
		errs["location_name"] = "Location name is required."
	}
	if t.LocationCode == "" {
		errs["location_code"] = "Location code is required."
	}
	if t.StatusID < 1 {
		errs["status_id"] = "Choose a status."
	}
	if t.PriorityID < 1 {
		errs["priority_id"] = "Choose a priority."
	}
	if t.CategoryID < 1 {
		errs["category_id"] = "Choose a category."
	}
	return t, errs
}