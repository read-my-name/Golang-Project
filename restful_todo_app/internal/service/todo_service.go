package service

import (
    "context"
    "fmt"
    "time"

    "github.com/read-my-name/restful_todo_app/pkg/models"
    "github.com/google/uuid"
)

type TodoService struct {
    storage    Storage
    saveQueue  chan models.Todo
    errorChan  chan error
}

type Storage interface {
    Save(t models.Todo) error
    SaveAll([]models.Todo) error
    Load() ([]models.Todo, error)
}

func NewTodoService(storage Storage) *TodoService {
    svc := &TodoService{
        storage:    storage,
        saveQueue:  make(chan models.Todo, 100),
        errorChan:  make(chan error, 10),
    }
    
    go svc.startSaveWorker(context.Background())
    go svc.startErrorHandler(context.Background())
    
    return svc
}

func isValidStatus(s models.Status) bool {
    switch s {
        case models.StatusNotStarted, 
            models.StatusInProgress,
            models.StatusOnHold, 
            models.StatusCompleted, 
            models.StatusArchived:
            return true
    }
    return false
}

func isValidPriority(p models.Priority) bool {
    switch p {
        case models.PriorityLow, 
            models.PriorityMedium,
            models.PriorityHigh, 
            models.PriorityCritical:
            return true
    }
    return false
}

func isValidPeriod(p string) bool {
    switch p {
        case "all", 
            "today", 
            "week", 
            "month":
            return true
    }
    return false
}

func isValidTitle(t string) bool {
    if len(t) > 100 {
        fmt.Println("Title should not be longer than 100 characters")
        return false
    }
    return true
}

func isValidDescription(d string) bool {
    if len(d) > 200 {
        fmt.Println("Description should not be longer than 200 characters")
        return false
    }
    return true
}

func isValidDueDate(d time.Time) bool {
    if d.Before(time.Now()) {
        fmt.Println("Due date should not be in the past")
        return false
    }
    return true
}

func isValidLabel(l string) bool {
    if len(l) > 20 {
        fmt.Println("Label should not be longer than 20 characters")
        return false
    }
    return true
}

func isValidSubtask(s models.Subtask) bool {
    if len(s.Title) > 100 {
        fmt.Println("Subtask title should not be longer than 100 characters")
        return false
    }
    return true
}

func (s *TodoService) startSaveWorker(ctx context.Context) {
    for {
        select {
            case todo := <-s.saveQueue:
                fmt.Printf("Received todo with ID %s for saving", todo.ID)
                if err := s.storage.Save(todo); err != nil {
                    fmt.Printf("Error saving todo with ID %s: %v", todo.ID, err)
                    s.errorChan <- err
                } else {
                    fmt.Printf("Successfully saved todo with ID %s", todo.ID)
                }
            case <-ctx.Done():
                fmt.Println("Save worker context canceled, exiting")
                return
        }
    }
}

func (s *TodoService) startErrorHandler(ctx context.Context) {
    for {
        select {
            case err := <-s.errorChan:
                fmt.Printf("Error encountered: %v", err)
            case <-ctx.Done():
                fmt.Println("Error handler context canceled, exiting")
                return
        }
    }
}

func (s *TodoService) GetTodos(filter models.TodoFilter) ([]models.Todo, error) {
    todos, err := s.storage.Load()
    if err != nil {
        fmt.Printf("Error loading todos: %v", err)
        return nil, err
    }
    
    if filter.Period != "all" {
        return s.CategorizeTodos(todos, filter), nil
    }
    return todos, nil
}

func (s *TodoService) AddTodo(t models.Todo) error {
    t.ID = uuid.New().String()
    
    // Set timestamps
    now := time.Now().UTC()
    t.CreatedAt = now
    t.UpdatedAt = now
    
    // Validate status
    if !isValidStatus(t.Status) {
        return fmt.Errorf("invalid status: %s", t.Status)
    }
    
    // Validate priority
    if !isValidPriority(t.Priority) {
        return fmt.Errorf("invalid priority: %s", t.Priority)
    }

    if !isValidDueDate(t.DueDate) {
        return fmt.Errorf("invalid due date: %s", t.DueDate)
    }
    
    // Validate title
    if !isValidTitle(t.Title) {
        return fmt.Errorf("invalid title: %s", t.Title)
    }
    
    // Validate description
    if !isValidDescription(t.Description) {
        return fmt.Errorf("invalid description: %s", t.Description)
    }

    // Validate subtasks
    for _, subtask := range t.Subtasks {
        if !isValidSubtask(subtask) {
            return fmt.Errorf("invalid subtask: %s", subtask.Title)
        }
    }
    
    // Validate labels
    for _, label := range t.Labels {
        if !isValidLabel(label) {
            return fmt.Errorf("invalid label: %s", label)
        }
    }
    
    select {
        case s.saveQueue <- t:
            return nil
        default:
            return fmt.Errorf("save queue is full")
    }
}

func (s *TodoService) UpdateTodo(id string, updatedTodo models.Todo) error {
    existing, err := s.storage.Load()
    if err != nil {
        return err
    }
    
    updatedTodo.UpdatedAt = time.Now().UTC()
    
    // Validate status
    if !isValidStatus(updatedTodo.Status) {
        return fmt.Errorf("invalid status: %s", updatedTodo.Status)
    }
    
    // Validate priority
    if !isValidPriority(updatedTodo.Priority) {
        return fmt.Errorf("invalid priority: %s", updatedTodo.Priority)
    }
    
    // Validate due date
    if !isValidDueDate(updatedTodo.DueDate) {
        return fmt.Errorf("invalid due date: %s", updatedTodo.DueDate)
    }
    
    // Validate title
    if !isValidTitle(updatedTodo.Title) {
        return fmt.Errorf("invalid title: %s", updatedTodo.Title)
    }
    
    // Validate description
    if !isValidDescription(updatedTodo.Description) {
        return fmt.Errorf("invalid description: %s", updatedTodo.Description)
    }

    // Validate subtasks
    for _, subtask := range updatedTodo.Subtasks {
        if !isValidSubtask(subtask) {
            return fmt.Errorf("invalid subtask: %s", subtask.Title)
        }
    }

    // Validate labels
    for _, label := range updatedTodo.Labels {
        if !isValidLabel(label) {
            return fmt.Errorf("invalid label: %s", label)
        }
    }

    for i, todo := range existing {
        if todo.ID == id {
            existing[i] = updatedTodo
            return s.storage.Save(updatedTodo)
        }
    }
    return fmt.Errorf("todo with ID %s not found", id)
}

func (s *TodoService) DeleteTodo(id string) error {
    existing, err := s.storage.Load()
    if err != nil {
        return err
    }
    
    for i, todo := range existing {
        if todo.ID == id {
            existing = append(existing[:i], existing[i+1:]...)
            return s.storage.Save(existing[i])
        }
    }
    return fmt.Errorf("todo with ID %s not found", id)
}

func (s *TodoService) LoadInitialData() ([]models.Todo, error) {
    return s.storage.Load()
}

func (s *TodoService) CategorizeTodos(todos []models.Todo, filter models.TodoFilter) []models.Todo {
    now := time.Now()
    var filtered []models.Todo
    
    for _, t := range todos {
        switch filter.Period {
        case "today":
            if t.DueDate.Truncate(24*time.Hour).Equal(now.Truncate(24*time.Hour)) {
                filtered = append(filtered, t)
            }
        case "week":
            year, week := t.DueDate.ISOWeek()
            currentYear, currentWeek := now.ISOWeek()
            if year == currentYear && week == currentWeek {
                filtered = append(filtered, t)
            }
        case "month":
            if t.DueDate.Month() == now.Month() && t.DueDate.Year() == now.Year() {
                filtered = append(filtered, t)
            }
        }
    }
    return filtered
}