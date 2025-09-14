package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"ToDo-List/internal/adapters/logger"
	"ToDo-List/internal/core/domain"
	"ToDo-List/internal/core/ports"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type TodoHandler struct {
	todoService ports.ToDoService
	logger      *logger.Logger
}

func NewTodoHandler(todoService ports.ToDoService, logger *logger.Logger) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
		logger:      logger,
	}
}

// CreateTodoHandler - POST /api/todo
func (h *TodoHandler) CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received POST /api/todo request")
	
	var todo domain.ToDo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		h.logger.Warn("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	todo.Id = uuid.NewString()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	if strings.TrimSpace(todo.Todo) == "" {
		h.logger.Warn("Missing 'todo' field in request")
		http.Error(w, "Missing 'todo' field", http.StatusBadRequest)
		return
	}

	// Message и Deadline делаем опциональными
	if strings.TrimSpace(todo.Priority) == "" {
		todo.Priority = "medium"
	}
	if todo.Priority != "low" && todo.Priority != "medium" && todo.Priority != "high" {
		todo.Priority = "medium"
	}
	todo.Complete = false

	h.logger.Debug("Creating todo: %+v", todo)
	createdTodo, err := h.todoService.CreateTodo(r.Context(), todo)
	if err != nil {
		h.logger.Error("Failed to create todo: %v", err)
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Todo created successfully: %s", createdTodo.Id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTodo)
}

// GetTodosHandler - GET /api/todos
func (h *TodoHandler) GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Received GET /api/todos request")
	
	q := r.URL.Query()
	
	filter := ports.TodoFilter{
		Status:   q.Get("status"),
		OrderBy:  q.Get("orderBy"),
		OrderDir: q.Get("orderDir"),
		Period:   q.Get("period"),
	}
	
	// Валидация параметров
	if filter.Status != "" && filter.Status != "all" && filter.Status != "active" && 
	   filter.Status != "completed" && filter.Status != "overdue" {
		h.logger.Warn("Invalid status parameter: %s", filter.Status)
		http.Error(w, "Invalid status parameter", http.StatusBadRequest)
		return
	}
	
	if filter.OrderDir != "" && filter.OrderDir != "asc" && filter.OrderDir != "desc" {
		h.logger.Warn("Invalid orderDir parameter: %s", filter.OrderDir)
		http.Error(w, "Invalid orderDir parameter", http.StatusBadRequest)
		return
	}
	
	// Валидация OrderBy
	validOrderBy := map[string]bool{
		"created_at":   true,
		"deadline":     true,
		"priority":     true,
		"completed_at": true,
		"":             true, // пустое значение тоже валидно
	}
	
	if !validOrderBy[filter.OrderBy] {
		h.logger.Warn("Invalid orderBy parameter: %s", filter.OrderBy)
		http.Error(w, "Invalid orderBy parameter", http.StatusBadRequest)
		return
	}
	
	h.logger.Debug("Fetching todos with filter: %+v", filter)
	todos, err := h.todoService.GetAllTodosWithFilters(r.Context(), filter)
	if err != nil {
		h.logger.Error("Failed to get todos: %v", err)
		http.Error(w, "Failed to get todos", http.StatusInternalServerError)
		return
	}
	
	h.logger.Info("Returning %d todos", len(todos))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
// GetTodoByIdHandler - GET /api/todo/{id}
func (h *TodoHandler) GetTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Info("Received GET /api/todo/%s request", id)

	if id == "" {
		h.logger.Warn("Missing todo ID in request")
		http.Error(w, "Missing todo ID", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Fetching todo with ID: %s", id)
	todo, err := h.todoService.GetTodoById(r.Context(), id)
	if err != nil {
		h.logger.Error("Todo not found: %s, error: %v", id, err)
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	h.logger.Info("Todo found: %s", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodoByIdHandler - PUT /api/todo/{id}
func (h *TodoHandler) UpdateTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Info("Received PUT /api/todo/%s request", id)

	if id == "" {
		h.logger.Warn("Missing todo ID in request")
		http.Error(w, "Missing todo ID", http.StatusBadRequest)
		return
	}

	var todo domain.ToDo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		h.logger.Warn("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(todo.Todo) == "" {
		h.logger.Warn("Missing 'todo' field in update request")
		http.Error(w, "Missing 'todo' field", http.StatusBadRequest)
		return
	}

	todo.Id = id
	h.logger.Debug("Updating todo: %+v", todo)
	
	err = h.todoService.UpdateTodo(r.Context(), todo)
	if err != nil {
		h.logger.Error("Failed to update todo %s: %v", id, err)
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Todo updated successfully: %s", id)
	w.WriteHeader(http.StatusNoContent)
}

// DeleteTodoHandler - DELETE /api/todo/{id}
func (h *TodoHandler) DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Info("Received DELETE /api/todo/%s request", id)

	if id == "" {
		h.logger.Warn("Missing todo ID in request")
		http.Error(w, "Missing todo ID", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Deleting todo with ID: %s", id)
	err := h.todoService.DeleteTodo(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete todo %s: %v", id, err)
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Todo deleted successfully: %s", id)
	w.WriteHeader(http.StatusNoContent)
}

// CompleteTodoByIdHandler - POST /api/todo/complete/{id}
func (h *TodoHandler) CompleteTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	h.logger.Info("Received POST /api/todo/complete/%s request", id)

	if id == "" {
		h.logger.Warn("Missing todo ID in complete request")
		http.Error(w, "Missing todo ID", http.StatusBadRequest)
		return
	}

	h.logger.Debug("Completing todo with ID: %s", id)
	err := h.todoService.CompleteTodoById(r.Context(), id)
	if err != nil {
		h.logger.Error("Failed to complete todo %s: %v", id, err)
		http.Error(w, "Failed to complete todo", http.StatusInternalServerError)
		return
	}

	h.logger.Info("Todo completed successfully: %s", id)
	w.WriteHeader(http.StatusNoContent)
}