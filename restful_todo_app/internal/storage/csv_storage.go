package storage

import (
    "encoding/csv"
    // "fmt"
    "os"
	"strconv"
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

// func ensureFileExists(filePath string) error {
//     // Check if the file exists
//     if _, err := os.Stat(filePath); os.IsNotExist(err) {
//         // Create the file if it doesn't exist
//         file, err := os.Create(filePath)
//         if err != nil {
//             return fmt.Errorf("failed to create file: %w", err)
//         }
//         defer file.Close()
//         fmt.Println("File created:", filePath)
//     }
//     return nil
// }

func (s *CSVStorage) Save(t models.Todo) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    file, err := os.OpenFile(s.filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
    if err != nil {
        return err
    }
    defer file.Close()
    
    writer := csv.NewWriter(file)
    defer writer.Flush()
    
    return writer.Write([]string{
        strconv.Itoa(t.ID),
        t.Title,
        strconv.FormatBool(t.Completed),
        t.DueDate.Format(time.RFC3339),
        t.CreatedAt.Format(time.RFC3339),
    })
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
        id, _ := strconv.Atoi(record[0])
        completed, _ := strconv.ParseBool(record[2])
        dueDate, _ := time.Parse(time.RFC3339, record[3])
        createdAt, _ := time.Parse(time.RFC3339, record[4])
        
        todos = append(todos, models.Todo{
            ID:          id,
            Title:       record[1],
            Completed:   completed,
            DueDate:     dueDate,
            CreatedAt:   createdAt,
        })
    }
    return todos, nil
}