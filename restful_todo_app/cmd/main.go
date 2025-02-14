package main

import (
    "log"
    
    "github.com/gin-gonic/gin"
    "github.com/read-my-name/restful_todo_app/api/handlers"
    "github.com/read-my-name/restful_todo_app/internal/service"
    "github.com/read-my-name/restful_todo_app/internal/storage"
    // "github.com/read-my-name/restful_todo_app/pkg/models"
)

func main() {
    // Initialize dependencies
    csvStorage := storage.NewCSVStorage("D:/GitHub/Freelance/restful_todo_app/data/todos.csv")
    todoService := service.NewTodoService(csvStorage)
    todoHandler := handlers.NewTodoHandler(todoService)
    
    // Load existing todos
    if _, err := todoService.LoadInitialData(); err != nil {
        log.Fatalf("Failed to load initial data: %v", err)
    }
    
    router := gin.Default()
    
    // Routes
    router.GET("/todos", todoHandler.GetTodos)
    router.POST("/todos", todoHandler.AddTodo)
    router.PUT("/todos/:id", todoHandler.UpdateTodo)
    router.DELETE("/todos/:id", todoHandler.DeleteTodo)
    
    // Time-based categorization
    router.GET("/todos/today", todoHandler.GetTodayTodos)
    // router.GET("/todos/week", todoHandler.GetWeekTodos)
    // router.GET("/todos/month", todoHandler.GetMonthTodos)
    
    if err := router.Run(":8080"); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}