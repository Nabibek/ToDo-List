package ports

import (
	"context"

	"ToDo-List/internal/core/domain"
)

type PostgreRepo interface {
	GetAllTodos(ctx context.Context) ([]domain.ToDo, error)
	GetTodoById(ctx context.Context, id string) (domain.ToDo, error)
	DeleteTodoById(ctx context.Context, id string) error
	UpdateTodo(ctx context.Context, todo domain.ToDo) error
	CreateTodo(ctx context.Context, todo domain.ToDo) (domain.ToDo, error)
}
