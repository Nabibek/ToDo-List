package main

import (
	"context"
	"log"

	"todo-desktop/backend/repo"
	"todo-desktop/backend/services"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx    context.Context
	logger *log.Logger
	todoService *services.TodoService
}

func NewApp() *App {
	return &App{
		logger: log.New(log.Writer(), "TODO ", log.LstdFlags),
	}
}

// Startup - вызывается при запуске приложения
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	
	// Инициализация репозитория (SQLite для desktop)
	dbRepo := repo.NewSQLiteRepo() 
	a.todoService = services.NewTodoService(dbRepo, a.logger)
	
	a.logger.Println("Application started")
}

// Methods exposed to frontend
func (a *App) GetTodos() ([]models.Todo, error) {
	return a.todoService.GetAllTodos(a.ctx)
}

func (a *App) CreateTodo(todo models.Todo) error {
	return a.todoService.CreateTodo(a.ctx, todo)
}

func (a *App) UpdateTodo(todo models.Todo) error {
	return a.todoService.UpdateTodo(a.ctx, todo)
}

func (a *App) DeleteTodo(id string) error {
	return a.todoService.DeleteTodo(a.ctx, id)
}

// Desktop-specific functionality
func (a *App) ShowNotification(title, message string) {
	runtime.Notification(a.ctx, &runtime.NotificationOptions{
		Title:   title,
		Message: message,
	})
}
