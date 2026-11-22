package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

func TestInvalidTicketSubmissionReturns422(t *testing.T) {
	form := url.Values{
		"title":       {strings.Repeat("x", 81)},
		"location":    {"Room 12"},
		"category":    {"Not a category"},
		"description": {"A fictional test description."},
	}
	request := httptest.NewRequest(http.MethodPost, "/tickets/create", strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response := httptest.NewRecorder()

	(&App{}).ticketCreatePost(response, request)

	if response.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusUnprocessableEntity)
	}
	if !strings.Contains(response.Body.String(), "Use at most 80 code points.") {
		t.Error("response does not contain the title validation message")
	}
	if !strings.Contains(response.Body.String(), "Choose a listed option.") {
		t.Error("response does not contain the category validation message")
	}
	if !strings.Contains(response.Body.String(), `aria-invalid="true"`) {
		t.Error("invalid fields are not marked in the response")
	}
}

func TestTicketIDRejectsMalformedAndNonpositiveValues(t *testing.T) {
	for _, value := range []string{"abc", "0", "-4"} {
		t.Run(value, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/tickets/"+value, nil)
			request.SetPathValue("id", value)
			response := httptest.NewRecorder()

			if _, ok := ticketID(response, request); ok {
				t.Fatal("ticketID accepted an invalid ID")
			}
			if response.Code != http.StatusNotFound {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusNotFound)
			}
		})
	}
}

func TestRouterRegistersTicketRoutes(t *testing.T) {
	router := newRouter(&App{})
	cases := []struct {
		method string
		path   string
		status int
	}{
		{http.MethodGet, "/", http.StatusSeeOther},
		{http.MethodGet, "/tickets/create", http.StatusOK},
		{http.MethodGet, "/tickets/not-an-id", http.StatusNotFound},
		{http.MethodGet, "/tickets/not-an-id/edit", http.StatusNotFound},
		{http.MethodPost, "/tickets/not-an-id/edit", http.StatusNotFound},
		{http.MethodGet, "/tickets/not-an-id/delete", http.StatusNotFound},
		{http.MethodPost, "/tickets/not-an-id/delete", http.StatusNotFound},
	}

	for _, test := range cases {
		t.Run(test.method+" "+test.path, func(t *testing.T) {
			request := httptest.NewRequest(test.method, test.path, nil)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != test.status {
				t.Fatalf("status = %d, want %d", response.Code, test.status)
			}
		})
	}
}

func TestTicketTemplatesRender(t *testing.T) {
	ticket := Ticket{
		ID: 1, Title: "Wi-Fi issue", Location: "Library",
		Category: "Network", Description: "Cannot connect.",
		Status: "New", StatusClass: statusClass("New"), Priority: "Medium",
		CreatedAt: time.Date(2026, 1, 2, 15, 4, 0, 0, time.UTC),
	}
	secondTicket := Ticket{
		ID: 2, Title: "Printer issue", Location: "Lab 4",
		Category: "Hardware", Status: "Resolved", StatusClass: statusClass("Resolved"),
	}
	cases := []struct {
		name string
		data any
	}{
		{"ticket-form.tmpl", FormData{}},
		{"ticket-form.tmpl", FormData{ID: 1, Title: ticket.Title}},
		{"tickets.tmpl", TicketListData{Items: []Ticket{ticket, secondTicket}}},
		{"ticket-detail.tmpl", ticket},
		{"ticket-delete.tmpl", ticket},
		{"software-form.tmpl", FormData{}},
		{"software-requests.tmpl", SoftwareRequestListData{}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			renderTemplate(response, http.StatusOK, test.name, test.data)
			if response.Code != http.StatusOK {
				t.Fatalf("render status = %d, want %d", response.Code, http.StatusOK)
			}
			if !strings.Contains(response.Body.String(), "<!doctype html>") {
				t.Fatal("response does not contain a rendered HTML document")
			}
		})
	}

	response := httptest.NewRecorder()
	renderTemplate(response, http.StatusOK, "tickets.tmpl", TicketListData{Items: []Ticket{ticket, secondTicket}})
	body := response.Body.String()
	for _, expected := range []string{
		`<a class="ticket-row" href="/tickets/1">`,
		`<a class="ticket-row" href="/tickets/2">`,
		`pill--new`,
		`pill--resolved`,
		ticket.Title,
		secondTicket.Title,
	} {
		if !strings.Contains(body, expected) {
			t.Errorf("ticket list does not contain %q", expected)
		}
	}
	if count := strings.Count(body, `class="ticket-row"`); count != 2 {
		t.Errorf("rendered ticket rows = %d, want 2", count)
	}

	emptyResponse := httptest.NewRecorder()
	renderTemplate(emptyResponse, http.StatusOK, "tickets.tmpl", TicketListData{})
	if emptyResponse.Code != http.StatusOK || !strings.Contains(emptyResponse.Body.String(), "No tickets yet") {
		t.Error("empty ticket list does not render a successful empty state")
	}
}

func TestTicketDetailsEscapeUserHTML(t *testing.T) {
	ticket := Ticket{
		ID: 7, Title: `<script>alert("x")</script>`, Location: "Library",
		Category: "Other", Description: "Description",
		Status: "New", StatusClass: statusClass("New"), Priority: "Low",
	}
	response := httptest.NewRecorder()
	renderTemplate(response, http.StatusOK, "ticket-detail.tmpl", ticket)
	body := response.Body.String()

	if strings.Contains(body, `<script>alert`) {
		t.Fatal("ticket title was rendered as executable HTML")
	}
	if !strings.Contains(body, `&lt;script&gt;alert(&#34;x&#34;)&lt;/script&gt;`) {
		t.Fatal("ticket title was not rendered as escaped text")
	}
}
