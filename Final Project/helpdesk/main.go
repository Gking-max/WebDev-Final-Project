package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	app := &App{Model: Model{DB: db}}
	server := &http.Server{
		Addr:              ":4000",
		Handler:           newRouter(app),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Print("Open http://localhost:4000/tickets")
	log.Fatal(server.ListenAndServe())
}

func newRouter(app *App) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/tickets", http.StatusSeeOther)
	})
	mux.HandleFunc("GET /tickets", app.listTickets)
	mux.HandleFunc("GET /tickets/create", app.ticketCreate)
	mux.HandleFunc("POST /tickets/create", app.ticketCreatePost)
	mux.HandleFunc("GET /tickets/{id}", app.ticketShow)
	mux.HandleFunc("GET /tickets/{id}/edit", app.ticketEdit)
	mux.HandleFunc("POST /tickets/{id}/edit", app.ticketEditPost)
	mux.HandleFunc("GET /tickets/{id}/delete", app.ticketDelete)
	mux.HandleFunc("POST /tickets/{id}/delete", app.ticketDeletePost)

	mux.HandleFunc("GET /software-requests/create", app.softwareCreate)
	mux.HandleFunc("POST /software-requests/create", app.softwareCreatePost)
	mux.HandleFunc("GET /software-requests", app.listSoftwareRequests)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	return mux
}
