package models

import (
    "time"
)

type Todo struct {
    ID          int       `json:"id"` // can change to uuid
    Title       string    `json:"title"`
    Description string    `json:"description"` // Added field
    Completed   bool      `json:"completed"`
    DueDate     time.Time `json:"due_date"`
    CreatedAt   time.Time `json:"created_at"`
}

type TodoFilter struct {
    Period string // today, week, month
}