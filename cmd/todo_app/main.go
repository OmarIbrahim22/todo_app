package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/OmarIbrahim22/todo_app/internal/adapters/http/handler"
	"github.com/OmarIbrahim22/todo_app/internal/adapters/store"
)

func main() {
	// open SQLite DB
	db, err := sql.Open("sqlite", "./todo.db")
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	// apply migrations
	if err := store.Migrate(db); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// create repo and router
	repo := store.NewSQLiteRepository(db)
	h := handler.New(repo)

	r := chi.NewRouter()
	r.Get("/healthz", h.Healthz)
	r.Get("/items", h.ListItems)
	r.Post("/items", h.CreateItem)
	r.Patch("/items/{id}/toggle", h.ToggleItem)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("Listening on :8080")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
