package ports

import (
	"context"

	"ToDo-List/internal/core/domain"
)

type ToDoService interface {
	CreateTodo(ctx context.Context, todo domain.ToDo) (domain.ToDo, error)
	GetTodoById(ctx context.Context, id string) (domain.ToDo, error)
	UpdateTodo(ctx context.Context, todo domain.ToDo) error
	DeleteTodo(ctx context.Context, id string) error
	CompleteTodoById(ctx context.Context, id string) error
	GetAllTodosWithFilters(ctx context.Context, filter TodoFilter) ([]domain.ToDo, error)
}
