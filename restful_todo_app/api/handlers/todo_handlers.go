package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/read-my-name/restful_todo_app/internal/service"
	"github.com/read-my-name/restful_todo_app/pkg/models"
)

type TodoHandler struct {
    service *service.TodoService
}

func NewTodoHandler(svc *service.TodoService) *TodoHandler {
    return &TodoHandler{service: svc}
}

func (h *TodoHandler) GetTodos(c *gin.Context) {
    filter := models.TodoFilter{
        Period: c.DefaultQuery("period", "all"),
    }
    
    todos, err := h.service.GetTodos(filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) AddTodo(c *gin.Context) {
    var newTodo models.Todo
    if err := c.BindJSON(&newTodo); err != nil {
        fmt.Printf("Error binding JSON: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    newTodo.CreatedAt = time.Now()
    fmt.Printf("Adding new todo: %+v\n", newTodo)
    if err := h.service.AddTodo(newTodo); err != nil {
        fmt.Printf("Error adding todo: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    fmt.Printf("Successfully added todo: %+v\n", newTodo)
    c.JSON(http.StatusCreated, newTodo)
}

func (h *TodoHandler) UpdateTodo(c *gin.Context){
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }

    var updatedTodo models.Todo
    if err := c.BindJSON(&updatedTodo); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if err := h.service.UpdateTodo(id, updatedTodo); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) 
        return
    }

    c.JSON(http.StatusOK, updatedTodo)
}

func (h *TodoHandler) DeleteTodo(c *gin.Context) {
    idParam := c.Param("id")
    id, err := strconv.Atoi(idParam)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
        return
    }    

    if err := h.service.DeleteTodo(id); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}

func (h *TodoHandler) GetTodayTodos(c *gin.Context) {
    filter := models.TodoFilter{
        Period: "today",
    }
    
    todos, err := h.service.GetTodos(filter)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, todos)
}
