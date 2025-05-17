package core

import "context"

// Repository defines persistence operations for Items.
type Repository interface {
    Create(ctx context.Context, item Item) error
    List(ctx context.Context, week int) ([]Item, error)
    ToggleDone(ctx context.Context, id string) error
}
