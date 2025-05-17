package handler

import (
    "context"
    "encoding/json"
    "net/http"
    "strconv"

    "github.com/go-chi/chi/v5"
    "github.com/OmarIbrahim22/todo_app/internal/core"
)

// Handler wraps core.Repository to serve HTTP endpoints.
type Handler struct {
    repo core.Repository
}

func New(repo core.Repository) *Handler {
    return &Handler{repo: repo}
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("ok"))
}

func (h *Handler) ListItems(w http.ResponseWriter, r *http.Request) {
    q := r.URL.Query().Get("week")
    week, _ := strconv.Atoi(q)
    items, err := h.repo.List(r.Context(), week)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    json.NewEncoder(w).Encode(items)
}

func (h *Handler) CreateItem(w http.ResponseWriter, r *http.Request) {
    var req struct {
        Description string `json:"description"`
        Priority    int    `json:"priority"`
        Week        int    `json:"week"`
    }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    item := core.Item{
        ID:          core.NewItem(req.Description, req.Priority).ID,
        Description: req.Description,
        Done:        false,
        Priority:    req.Priority,
        Week:        req.Week,
        CreatedAt:   time.Now(),
    }
    if err := h.repo.Create(r.Context(), item); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(item)
}

func (h *Handler) ToggleItem(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    if err := h.repo.ToggleDone(r.Context(), id); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusNoContent)
}
