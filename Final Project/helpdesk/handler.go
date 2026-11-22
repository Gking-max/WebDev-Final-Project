package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"
)

type App struct {
	Model Model
}

type FormData struct {
	ID                                                         int64
	Title, Location, Category, Description                     string
	TitleError, LocationError, CategoryError, DescriptionError string
}

type TicketListData struct {
	Items []Ticket
}

type SoftwareRequestListData struct {
	Items []SoftwareRequest
}

func (a *App) ticketCreate(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, http.StatusOK, "ticket-form.tmpl", FormData{})
}

func (a *App) ticketCreatePost(w http.ResponseWriter, r *http.Request) {
	form, ok := readAndValidateForm(w, r, []string{"Network", "Email", "Hardware", "Other"}, "Enter a title.")
	if !ok {
		if form.hasErrors() {
			renderTemplate(w, http.StatusUnprocessableEntity, "ticket-form.tmpl", form)
		}
		return
	}

	id, err := a.Model.InsertTicket(r.Context(), form.Title, form.Location, form.Category, form.Description)
	if err != nil {
		log.Printf("insert ticket: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/tickets/%d", id), http.StatusSeeOther)
}

func (a *App) listTickets(w http.ResponseWriter, r *http.Request) {
	tickets, err := a.Model.ListTickets(r.Context())
	if err != nil {
		log.Printf("list tickets: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, http.StatusOK, "tickets.tmpl", TicketListData{Items: tickets})
}

func (a *App) ticketShow(w http.ResponseWriter, r *http.Request) {
	id, ok := ticketID(w, r)
	if !ok {
		return
	}
	ticket, err := a.Model.GetTicket(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Printf("get ticket %d: %v", id, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, http.StatusOK, "ticket-detail.tmpl", ticket)
}

func (a *App) ticketEdit(w http.ResponseWriter, r *http.Request) {
	id, ok := ticketID(w, r)
	if !ok {
		return
	}
	ticket, err := a.Model.GetTicket(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Printf("get ticket %d for edit: %v", id, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, http.StatusOK, "ticket-form.tmpl", FormData{
		ID: ticket.ID, Title: ticket.Title, Location: ticket.Location,
		Category: ticket.Category, Description: ticket.Description,
	})
}

func (a *App) ticketEditPost(w http.ResponseWriter, r *http.Request) {
	id, ok := ticketID(w, r)
	if !ok {
		return
	}
	if _, err := a.Model.GetTicket(r.Context(), id); errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	} else if err != nil {
		log.Printf("get ticket %d before edit: %v", id, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	form, valid := readAndValidateForm(w, r, []string{"Network", "Email", "Hardware", "Other"}, "Enter a title.")
	if !valid {
		if form.hasErrors() {
			form.ID = id
			renderTemplate(w, http.StatusUnprocessableEntity, "ticket-form.tmpl", form)
		}
		return
	}

	rows, err := a.Model.UpdateTicket(r.Context(), id, form.Title, form.Location, form.Category, form.Description)
	if err != nil {
		log.Printf("update ticket %d: %v", id, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/tickets/%d", id), http.StatusSeeOther)
}

func (a *App) ticketDelete(w http.ResponseWriter, r *http.Request) {
	id, ok := ticketID(w, r)
	if !ok {
		return
	}
	ticket, err := a.Model.GetTicket(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		log.Printf("get ticket %d for delete: %v", id, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, http.StatusOK, "ticket-delete.tmpl", ticket)
}

func (a *App) ticketDeletePost(w http.ResponseWriter, r *http.Request) {
	id, ok := ticketID(w, r)
	if !ok {
		return
	}
	rows, err := a.Model.DeleteTicket(r.Context(), id)
	if err != nil {
		log.Printf("delete ticket %d: %v", id, err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	if rows == 0 {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, "/tickets", http.StatusSeeOther)
}

func (a *App) softwareCreate(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, http.StatusOK, "software-form.tmpl", FormData{})
}

func (a *App) softwareCreatePost(w http.ResponseWriter, r *http.Request) {
	form, ok := readAndValidateForm(w, r, []string{"Installation", "Upgrade", "Licence", "Other"}, "Enter a software name.")
	if !ok {
		if form.hasErrors() {
			renderTemplate(w, http.StatusUnprocessableEntity, "software-form.tmpl", form)
		}
		return
	}

	_, err := a.Model.InsertSoftwareRequest(r.Context(), form.Title, form.Location, form.Category, form.Description)
	if err != nil {
		log.Printf("insert software request: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/software-requests", http.StatusSeeOther)
}

func (a *App) listSoftwareRequests(w http.ResponseWriter, r *http.Request) {
	requests, err := a.Model.ListSoftwareRequests(r.Context())
	if err != nil {
		log.Printf("list software requests: %v", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	renderTemplate(w, http.StatusOK, "software-requests.tmpl", SoftwareRequestListData{Items: requests})
}

func readAndValidateForm(w http.ResponseWriter, r *http.Request, choices []string, titlePrompt string) (FormData, bool) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return FormData{}, false
	}
	form := FormData{
		Title:       r.PostForm.Get("title"),
		Location:    r.PostForm.Get("location"),
		Category:    r.PostForm.Get("category"),
		Description: r.PostForm.Get("description"),
	}
	form.Title = strings.TrimSpace(form.Title)
	form.Location = strings.TrimSpace(form.Location)
	form.Description = strings.TrimSpace(form.Description)

	if form.Title == "" {
		form.TitleError = titlePrompt
	} else if utf8.RuneCountInString(form.Title) > 80 {
		form.TitleError = "Use at most 80 code points."
	}
	if form.Location == "" {
		form.LocationError = "Enter a location."
	} else if utf8.RuneCountInString(form.Location) > 40 {
		form.LocationError = "Use at most 40 code points."
	}
	if !contains(choices, form.Category) {
		form.CategoryError = "Choose a listed option."
	}
	if form.Description == "" {
		form.DescriptionError = "Enter a description."
	} else if utf8.RuneCountInString(form.Description) > 500 {
		form.DescriptionError = "Use at most 500 code points."
	}
	return form, !form.hasErrors()
}

func (f FormData) hasErrors() bool {
	return f.TitleError != "" || f.LocationError != "" ||
		f.CategoryError != "" || f.DescriptionError != ""
}

func contains(choices []string, value string) bool {
	for _, choice := range choices {
		if value == choice {
			return true
		}
	}
	return false
}

func ticketID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return 0, false
	}
	return id, true
}
