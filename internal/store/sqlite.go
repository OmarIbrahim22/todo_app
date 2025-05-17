package store

import (
    "context"
    "database/sql"
    "embed"
    "fmt"
    "io/fs"
    "strconv"
    "time"

    "github.com/OmarIbrahim22/todo_app/internal/core"
    "modernc.org/sqlite"
)

//go:embed ../../../migrations/*.sql
var migrationsFS embed.FS

// Migrate applies all .sql files in the migrations directory in lexical order.
func Migrate(db *sql.DB) error {
    entries, err := fs.ReadDir(migrationsFS, ".")
    if err != nil {
        return err
    }
    for _, e := range entries {
        if e.IsDir() || filepath.Ext(e.Name()) != ".sql" {
            continue
        }
        b, err := migrationsFS.ReadFile(e.Name())
        if err != nil {
            return err
        }
        if _, err := db.Exec(string(b)); err != nil {
            return fmt.Errorf("exec %s: %w", e.Name(), err)
        }
    }
    return nil
}

type sqliteRepo struct {
    db *sql.DB
}

// NewSQLiteRepository constructs a Core.Repository backed by SQLite.
func NewSQLiteRepository(db *sql.DB) core.Repository {
    return &sqliteRepo{db: db}
}

func (s *sqliteRepo) Create(ctx context.Context, item core.Item) error {
    _, err := s.db.ExecContext(ctx, `
        INSERT INTO items(id, description, done, priority, week, created_at)
        VALUES(?, ?, ?, ?, ?, ?)
    `, item.ID.String(), item.Description, item.Done, item.Priority, item.Week, item.CreatedAt)
    return err
}

func (s *sqliteRepo) List(ctx context.Context, week int) ([]core.Item, error) {
    rows, err := s.db.QueryContext(ctx, `
        SELECT id, description, done, priority, week, created_at
        FROM items
        WHERE week = ?
        ORDER BY priority, created_at
    `, week)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []core.Item
    for rows.Next() {
        var it core.Item
        var idStr string
        if err := rows.Scan(&idStr, &it.Description, &it.Done, &it.Priority, &it.Week, &it.CreatedAt); err != nil {
            return nil, err
        }
        it.ID, _ = uuid.Parse(idStr)
        items = append(items, it)
    }
    return items, rows.Err()
}

func (s *sqliteRepo) ToggleDone(ctx context.Context, id string) error {
    _, err := s.db.ExecContext(ctx, `
        UPDATE items SET done = NOT done WHERE id = ?
    `, id)
    return err
}
