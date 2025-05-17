package core

import (
    "time"

    "github.com/google/uuid"
)

// Item represents a todo checklist item.
type Item struct {
    ID          uuid.UUID `json:"id"`
    Description string    `json:"description"`
    Done        bool      `json:"done"`
    Priority    int       `json:"priority"` // 1 (highest) through 5 (lowest)
    Week        int       `json:"week"`     // ISO week number
    CreatedAt   time.Time `json:"created_at"`
}

// NewItem constructs a new Item for the current week.
func NewItem(desc string, priority int) Item {
    now := time.Now()
    year, week := now.ISOWeek()
    return Item{
        ID:          uuid.New(),
        Description: desc,
        Done:        false,
        Priority:    priority,
        Week:        week + year*100, // encode year/week into one int
        CreatedAt:   now,
    }
}
