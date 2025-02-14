package storage

import (
    "encoding/csv"
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
        return err
    }

    file, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

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
    })
}

func (s *CSVStorage) SaveAll(todos []models.Todo) error {
    s.mu.Lock()
    defer s.mu.Unlock()

    dir := filepath.Dir(s.filePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }

    file, err := os.Create(s.filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    for _, t := range todos {
        err := writer.Write([]string{
            t.ID,
            t.Title,
            t.Description,
            string(t.Status),
            string(t.Priority),
            t.DueDate.Format(time.RFC3339),
            t.CreatedAt.Format(time.RFC3339),
            t.UpdatedAt.Format(time.RFC3339),
            strings.Join(t.Labels, "|"),
        })
        if err != nil {
            return err
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
        return nil, err
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        return nil, err
    }

    var todos []models.Todo
    for _, record := range records {
        if len(record) < 9 {
            continue // Skip invalid records
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
            Labels:      strings.Split(record[8], "|"),
        }

        // Remove empty labels
        var cleanLabels []string
        for _, label := range todo.Labels {
            if label != "" {
                cleanLabels = append(cleanLabels, label)
            }
        }
        todo.Labels = cleanLabels

        todos = append(todos, todo)
    }

    return todos, nil
}