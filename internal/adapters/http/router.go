package http

import (
	"net/http"
	"path/filepath"
	"strings"

	"ToDo-List/internal/adapters/http/handlers"
	"ToDo-List/internal/adapters/logger"
	"ToDo-List/internal/application/service"
	"ToDo-List/internal/core/ports"

	"github.com/gorilla/mux"
)

func NewRouter(repo ports.PostgreRepo, appLogger *logger.Logger) *mux.Router {
	router := mux.NewRouter()

	appLogger.Info("Initializing HTTP router...")

	todoService := service.NewToDoService(repo, appLogger) // передаем логгер в сервис
	todoHandler := handlers.NewTodoHandler(todoService, appLogger)

	// Создаем подроутер для API с префиксом /api
	apiRouter := router.PathPrefix("/api").Subrouter()

	// POST /api/todo
	apiRouter.HandleFunc("/todo", todoHandler.CreateTodoHandler).Methods(http.MethodPost)

	// GET /api/todos
	apiRouter.HandleFunc("/todos", todoHandler.GetTodosHandler).Methods(http.MethodGet)

	// Health check
	router.HandleFunc("/health", healthHandler(repo, appLogger)).Methods(http.MethodGet)

	// POST /api/todo/complete/{id}
	apiRouter.HandleFunc("/todo/complete/{id}", todoHandler.CompleteTodoByIdHandler).Methods(http.MethodPost)

	// GET /api/todo/{id}
	apiRouter.HandleFunc("/todo/{id}", todoHandler.GetTodoByIdHandler).Methods(http.MethodGet)
	// PUT /api/todo/{id}
	apiRouter.HandleFunc("/todo/{id}", todoHandler.UpdateTodoByIdHandler).Methods(http.MethodPut)
	// DELETE /api/todo/{id}
	apiRouter.HandleFunc("/todo/{id}", todoHandler.DeleteTodoHandler).Methods(http.MethodDelete)

	// Обслуживание статических файлов фронтенда
	router.PathPrefix("/").Handler(customFileServer("./web", appLogger))

	appLogger.Info("HTTP router initialized successfully")
	return router
}

func customFileServer(root string, appLogger *logger.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		appLogger.Debug("Serving static file: %s", r.URL.Path)
		
		switch {
		case strings.HasSuffix(r.URL.Path, ".js"):
			w.Header().Set("Content-Type", "application/javascript")
		case strings.HasSuffix(r.URL.Path, ".css"):
			w.Header().Set("Content-Type", "text/css")
		case strings.HasSuffix(r.URL.Path, ".html"):
			w.Header().Set("Content-Type", "text/html")
		}
		
		http.ServeFile(w, r, filepath.Join(root, r.URL.Path))
	})
}

func healthHandler(repo ports.PostgreRepo, appLogger *logger.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		appLogger.Debug("Health check request from %s", r.RemoteAddr)
		
		if err := repo.Ping(); err != nil {
			appLogger.Error("Database health check failed: %v", err)
			http.Error(w, "Database connection error", http.StatusInternalServerError)
			return
		}
		
		appLogger.Debug("Health check passed")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}