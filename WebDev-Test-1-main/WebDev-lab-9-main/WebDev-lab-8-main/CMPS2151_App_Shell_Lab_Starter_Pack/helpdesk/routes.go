package main

import "net/http"

func (a *App) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", a.home)
	mux.HandleFunc("GET /tickets", a.listTickets)
	mux.HandleFunc("GET /tickets/create", a.createForm)
	mux.HandleFunc("POST /tickets/create", a.createSubmit)
	mux.HandleFunc("GET /tickets/{id}", a.showTicket)
	mux.HandleFunc("GET /tickets/{id}/edit", a.editForm)
	mux.HandleFunc("POST /tickets/{id}/edit", a.editSubmit)
	mux.HandleFunc("GET /tickets/{id}/delete", a.deleteConfirm)
	mux.HandleFunc("POST /tickets/{id}/delete", a.deleteSubmit)

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("GET /css/", http.StripPrefix("/css/", fs))
	return mux
}