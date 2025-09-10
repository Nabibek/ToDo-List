package http

import (
	"net/http"

	"ToDo-List/internal/adapters/http/handlers"
	"ToDo-List/internal/application/service"
	"ToDo-List/internal/core/ports"
)

func NewRouter(repo ports.PostgreRepo) *http.ServeMux {
	mux := http.NewServeMux()

	todoService := service.NewToDoService(repo)
	todoHandler := handlers.NewTodoHandler(todoService)

	// POST /todo
	mux.HandleFunc("/todo", todoHandler.CreateTodoHandler)

	// GET /todos
	mux.HandleFunc("/todos", todoHandler.GetTodosHandler)

	// GET /todo/{id}, PUT /todo/{id}, DELETE /todo/{id}
	mux.HandleFunc("/todo/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoHandler.GetTodoByIdHandler(w, r)
		case http.MethodPut:
			todoHandler.UpdateTodoByIdHandler(w, r)
		case http.MethodDelete:
			todoHandler.DeleteTodoHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	return mux
}
