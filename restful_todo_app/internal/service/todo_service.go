package service

import (
    "context"
    "fmt"
    "time"
    "github.com/read-my-name/restful_todo_app/pkg/models"
)

type TodoService struct {
    storage    Storage
    saveQueue  chan models.Todo
    errorChan  chan error
}

type Storage interface {
    Save(t models.Todo) error
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

func (s *TodoService) startSaveWorker(ctx context.Context) {
    for {
        select {
            case todo := <-s.saveQueue:
                fmt.Printf("Received todo with ID %d for saving", todo.ID)
                if err := s.storage.Save(todo); err != nil {
                    fmt.Printf("Error saving todo with ID %d: %v", todo.ID, err)
                    s.errorChan <- err
                } else {
                    fmt.Printf("Successfully saved todo with ID %d", todo.ID)
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
        return nil, err
    }
    
    if filter.Period != "all" {
        return s.CategorizeTodos(todos, filter), nil
    }
    return todos, nil
}

func (s *TodoService) AddTodo(t models.Todo) error {
    existing, _ := s.storage.Load()
    maxID := 0
    for _, todo := range existing {
        if todo.ID > maxID {
            maxID = todo.ID
        }
    }
    t.ID = maxID + 1
    
    select {
        case s.saveQueue <- t:
            return nil
        default:
            return fmt.Errorf("save queue is full")
    }
}

func (s *TodoService) UpdateTodo(id int, updatedTodo models.Todo) error {
    existing, err := s.storage.Load()
    if err != nil {
        return err
    }
    
    for i, todo := range existing {
        if todo.ID == id {
            existing[i] = updatedTodo
            return s.storage.Save(updatedTodo)
        }
    }
    return fmt.Errorf("todo with ID %d not found", id)
}

func (s *TodoService) DeleteTodo(id int) error {
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
    return fmt.Errorf("todo with ID %d not found", id)
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