package storage

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "strings"
    "sync"
    "time"
    
    "github.com/read-my-name/restful_todo_app/pkg/models"
)

type CSVStorage struct {
    filePath string
    mu       sync.RWMutex
}

func NewCSVStorage(filePath string) *CSVStorage {
    return &CSVStorage{filePath: filePath}
}

func (s *CSVStorage) Save(t models.Todo) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Create directory if not exists
    dir := filepath.Dir(s.filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    file, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        return fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Serialize subtasks to JSON
    subtasksJSON, err := json.Marshal(t.Subtasks)
    if err != nil {
        return fmt.Errorf("failed to marshal subtasks: %w", err)
    }

    return writer.Write([]string{
        t.ID,
        t.Title,
        t.Description,
        string(t.Status),
        string(t.Priority),
        t.DueDate.Format(time.RFC3339),
        t.CreatedAt.Format(time.RFC3339),
        t.UpdatedAt.Format(time.RFC3339),
        strings.Join(t.Labels, "|"),
        string(subtasksJSON),
    })
}

func (s *CSVStorage) SaveAll(todos []models.Todo) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    dir := filepath.Dir(s.filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory: %w", err)
    }

    file, err := os.Create(s.filePath)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, t := range todos {
        subtasksJSON, err := json.Marshal(t.Subtasks)
        if err != nil {
            return fmt.Errorf("failed to marshal subtasks: %w", err)
        }

        err = writer.Write([]string{
            t.ID,
            t.Title,
            t.Description,
            string(t.Status),
            string(t.Priority),
            t.DueDate.Format(time.RFC3339),
            t.CreatedAt.Format(time.RFC3339),
            t.UpdatedAt.Format(time.RFC3339),
            strings.Join(t.Labels, "|"),
            string(subtasksJSON),
        })
        if err != nil {
            return fmt.Errorf("failed to write record: %w", err)
        }
    }
    return nil
}

func (s *CSVStorage) Load() ([]models.Todo, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    file, err := os.Open(s.filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return []models.Todo{}, nil
        }
        return nil, fmt.Errorf("failed to open file: %w", err)
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, fmt.Errorf("failed to read CSV: %w", err)
    }

    var todos []models.Todo
    for idx, record := range records {
        // Handle both old (9 columns) and new (10 columns) formats
        if len(record) < 9 {
            return nil, fmt.Errorf("invalid record at line %d: expected at least 9 columns, got %d", idx+1, len(record))
        }

        dueDate, _ := time.Parse(time.RFC3339, record[5])
        createdAt, _ := time.Parse(time.RFC3339, record[6])
        updatedAt, _ := time.Parse(time.RFC3339, record[7])

        todo := models.Todo{
            ID:          record[0],
            Title:       record[1],
            Description: record[2],
            Status:      models.Status(record[3]),
            Priority:    models.Priority(record[4]),
            DueDate:     dueDate,
            CreatedAt:   createdAt,
            UpdatedAt:   updatedAt,
            Labels:      cleanStrings(strings.Split(record[8], "|")),
        }

        // Handle subtasks (column 9) if present
        if len(record) >= 10 {
            var subtasks []models.Subtask
            if record[9] != "" {
                if err := json.Unmarshal([]byte(record[9]), &subtasks); err != nil {
                    return nil, fmt.Errorf("failed to unmarshal subtasks at line %d: %w", idx+1, err)
                }
            }
            todo.Subtasks = subtasks
        }

        todos = append(todos, todo)
    }

    return todos, nil
}

// Helper function to clean empty strings from slices
func cleanStrings(slice []string) []string {
    var clean []string
    for _, s := range slice {
        if s != "" {
            clean = append(clean, s)
        }
    }
    return clean
}