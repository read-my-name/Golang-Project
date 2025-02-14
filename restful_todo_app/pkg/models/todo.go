// package models

// import (
//     "time"
// )

// type Todo struct {
//     ID          int       `json:"id"` // can change to uuid
//     Title       string    `json:"title"`
//     Description string    `json:"description"` // Added field
//     Completed   bool      `json:"completed"`
//     DueDate     time.Time `json:"due_date"`
//     CreatedAt   time.Time `json:"created_at"`
// }

// type TodoFilter struct {
//     Period string // today, week, month
// }

package models

import (
    "time"
)

type Status string

const (
    StatusNotStarted  Status = "not_started"
    StatusInProgress  Status = "in_progress"
    StatusOnHold      Status = "on_hold"
    StatusCompleted   Status = "completed"
    StatusArchived    Status = "archived"
)

type Priority string

const (
    PriorityLow      Priority = "low"
    PriorityMedium   Priority = "medium"
    PriorityHigh     Priority = "high"
    PriorityCritical Priority = "critical"
)

type Todo struct {
    ID          string     `json:"id"`          // Changed to UUID
    Title       string     `json:"title"`
    Description string     `json:"description"`
    Status      Status     `json:"status"`      // New status field
    Priority    Priority   `json:"priority"`    // New priority field
    DueDate     time.Time  `json:"due_date"`
    CreatedAt   time.Time  `json:"created_at"`
    UpdatedAt   time.Time  `json:"updated_at"`  // New field
    Labels      []string   `json:"labels"`      // New labels/tags
    Subtasks    []Subtask  `json:"subtasks"`    // New subtasks
}

type Subtask struct {
    ID          string    `json:"id"`
    Title       string    `json:"title"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type TodoFilter struct {
    Period      string      `json:"period"`   // today, week, month
    Statuses    []Status    `json:"statuses"` // Filter by multiple statuses
    Priorities  []Priority  `json:"priorities"`
    Labels      []string    `json:"labels"`
}