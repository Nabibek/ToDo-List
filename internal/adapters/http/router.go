package http

import (
	"net/http"

	"ToDo-List/internal/adapters/http/handlers"
	"ToDo-List/internal/application/service"
	"ToDo-List/internal/core/ports"

	"github.com/gorilla/mux"
)

func NewRouter(repo ports.PostgreRepo) *mux.Router {
	router := mux.NewRouter()

	todoService := service.NewToDoService(repo)
	todoHandler := handlers.NewTodoHandler(todoService)

	// POST /todo
	router.HandleFunc("/todo", todoHandler.CreateTodoHandler).Methods(http.MethodPost)

	// GET /todos
	router.HandleFunc("/todos", todoHandler.GetTodosHandler).Methods(http.MethodGet)

	// Health check
	router.HandleFunc("/health", healthHandler(repo)).Methods(http.MethodGet)

	// POST /todo/complete/{id}
	router.HandleFunc("/todo/complete/{id}", todoHandler.CompleteTodoByIdHandler).Methods(http.MethodPost)

	// GET /todo/{id}
	router.HandleFunc("/todo/{id}", todoHandler.GetTodoByIdHandler).Methods(http.MethodGet)
	// PUT /todo/{id}
	router.HandleFunc("/todo/{id}", todoHandler.UpdateTodoByIdHandler).Methods(http.MethodPut)
	// DELETE /todo/{id}
	router.HandleFunc("/todo/{id}", todoHandler.DeleteTodoHandler).Methods(http.MethodDelete)

	// router.HandleFunc("/todos/order", todoHandler.GetTodosOrderByHandler).Methods(http.MethodGet)
	// router.HandleFunc("/todos/period", todoHandler.GetTodosByPeriodHandler).Methods(http.MethodGet)
	// router.HandleFunc("/todos/status", todoHandler.GetTodosByStatusHandler).Methods(http.MethodGet)

	return router
}

func healthHandler(repo ports.PostgreRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := repo.Ping(); err != nil {
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
