package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"ToDo-List/internal/core/domain"
	"ToDo-List/internal/core/ports"
)

type TodoHandler struct {
	todoService ports.ToDoService
}

func NewTodoHandler(todoService ports.ToDoService) *TodoHandler {
	return &TodoHandler{todoService: todoService}
}

// CreateTodoHandler - POST /todo
func (h *TodoHandler) CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed suka?", http.StatusMethodNotAllowed)
		return
	}

	var todo domain.ToDo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if todo.Todo == "" {
		http.Error(w, "Missing 'todo' field", http.StatusBadRequest)
		return
	}
	if todo.Message == "" {
		http.Error(w, "Missing 'message' field", http.StatusBadRequest)
		return
	}
	if todo.Deadline.IsZero() {
		http.Error(w, "Missing or invalid 'deadline' field", http.StatusBadRequest)
		return
	}

	createdTodo, err := h.todoService.CreateTodo(r.Context(), todo)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdTodo)
}

// GetTodosHandler - GET /todos
func (h *TodoHandler) GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	todos, err := h.todoService.GetAllTodos(r.Context())
	if err != nil {
		http.Error(w, "Failed to get todos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// GetTodoByIdHandler - GET /todo/{id}
func (h *TodoHandler) GetTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/todo/")
	if id == "" {
		http.Error(w, "Missing todo ID", http.StatusBadRequest)
		return
	}

	todo, err := h.todoService.GetTodoById(r.Context(), id)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// UpdateTodoByIdHandler - PUT /todo/{id}
func (h *TodoHandler) UpdateTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPut {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/todo/")
	if id == "" {
		http.Error(w, "Missing todo ID", http.StatusBadRequest)
		return
	}

	var todo domain.ToDo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todo.Id = id

	err = h.todoService.UpdateTodo(r.Context(), todo)
	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}

// DeleteTodoHandler - DELETE /todo/{id}
func (h *TodoHandler) DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/todo/")
	if id == "" {
		http.Error(w, "Missing todo ID", http.StatusBadRequest)
		return
	}

	err := h.todoService.DeleteTodo(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TodoHandler) GetTodosByStatusHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	status := r.URL.Query().Get("status")
	if status == "" {
		http.Error(w, "Missing status parameter", http.StatusBadRequest)
		return
	}
	todos, err := h.todoService.GetTodosByStatus(r.Context(), status)
	if err != nil {
		http.Error(w, "Failed to get todos by status", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
func (h *TodoHandler) GetTodosByPeriodHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	if start == "" || end == "" {
		http.Error(w, "Missing start or end parameter", http.StatusBadRequest)
		return
	}
	todos, err := h.todoService.GetTodoByPeriod(r.Context(), start, end)
	if err != nil {
		http.Error(w, "Failed to get todos by period", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
func (h *TodoHandler) CompleteTodoByIdHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	id := strings.TrimPrefix(r.URL.Path, "/todo/complete/")
	if id == "" {
		http.Error(w, "Missing todo ID", http.StatusBadRequest)
		return
	}
	err := h.todoService.CompleteTodoById(r.Context(), id)
	if err != nil {
		http.Error(w, "Failed to complete todo", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *TodoHandler) GetTodosWithFilterHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received %s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	filters := make(map[string]string)
	queryParams := r.URL.Query()
	for key, values := range queryParams {
		if len(values) > 0 {
			filters[key] = strings.ToLower(values[0])
		}
	}
	todos, err := h.todoService.GetTodosWithFilter(r.Context(), filters)
	if err != nil {
		http.Error(w, "Failed to get todos with filters", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
