package ports

import (
	"context"

	"ToDo-List/internal/core/domain"
)

type ToDoService interface {
	CreateTodo(ctx context.Context, todo domain.ToDo) (domain.ToDo, error)
	GetAllTodos(ctx context.Context) ([]domain.ToDo, error)
	GetTodoById(ctx context.Context, id string) (domain.ToDo, error)
	UpdateTodo(ctx context.Context, todo domain.ToDo) error
	DeleteTodo(ctx context.Context, id string) error
	GetTodosByStatus(ctx context.Context, status string) ([]domain.ToDo, error)
	CompleteTodoById(ctx context.Context, id string) error
	GetTodoByPeriod(ctx context.Context, start string, end string) ([]domain.ToDo, error)
	GetTodosWithFilter(ctx context.Context, filters map[string]string) ([]domain.ToDo, error)
}
