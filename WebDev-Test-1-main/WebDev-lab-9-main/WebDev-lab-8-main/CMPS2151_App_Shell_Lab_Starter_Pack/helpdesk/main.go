package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	db, err := openDB()
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close()

	cache, err := loadTemplates("templates")
	if err != nil {
		log.Fatalf("templates: %v", err)
	}

	app := &App{Model: &Model{DB: db}, Cache: cache}

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":4000"
	}
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, app.routes()))
}

func loadTemplates(dir string) (map[string]*template.Template, error) {
	pages, err := filepath.Glob(filepath.Join(dir, "*.html"))
	if err != nil {
		return nil, err
	}
	cache := map[string]*template.Template{}
	for _, p := range pages {
		name := filepath.Base(p)
		if name == "base.html" {
			continue
		}
		ts, err := template.New(name).Funcs(template.FuncMap{"statusMod": statusMod}).ParseFiles(
			filepath.Join(dir, "base.html"), p,
		)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}