package ports

import (
	"context"

	"ToDo-List/internal/core/domain"
)

type TodoFilter struct {
	Status   string // "all", "active", "completed", "overdue"
	OrderBy  string // "created_at", "deadline", "priority", "completed_at"
	OrderDir string // "asc", "desc"
	Period   string // "today", "week", "month", "overdue"
}

type PostgreRepo interface {
	GetAllTodosWithFilters(ctx context.Context, filter TodoFilter) ([]domain.ToDo, error)
	GetTodoById(ctx context.Context, id string) (domain.ToDo, error)
	DeleteTodoById(ctx context.Context, id string) error
	UpdateTodo(ctx context.Context, todo domain.ToDo) error
	CreateTodo(ctx context.Context, todo domain.ToDo) (domain.ToDo, error)
	Ping() error
}