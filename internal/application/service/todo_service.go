package service

import (
	"context"
	"time"

	"ToDo-List/internal/core/domain"
	"ToDo-List/internal/core/ports"

	"github.com/google/uuid"
)

type TodoService struct {
	repo ports.PostgreRepo
}

func NewToDoService(repo ports.PostgreRepo) ports.ToDoService {
	return &TodoService{repo: repo}
}

func (s *TodoService) CreateTodo(ctx context.Context, todo domain.ToDo) (domain.ToDo, error) {
	todo.Id = uuid.NewString()
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()

	return s.repo.CreateTodo(ctx, todo)
}

func (s *TodoService) GetAllTodos(ctx context.Context) ([]domain.ToDo, error) {
	return s.repo.GetAllTodos(ctx)
}

func (s *TodoService) GetTodoById(ctx context.Context, id string) (domain.ToDo, error) {
	return s.repo.GetTodoById(ctx, id)
}

func (s *TodoService) UpdateTodo(ctx context.Context, todo domain.ToDo) error {
	todo.UpdatedAt = time.Now()
	return s.repo.UpdateTodo(ctx, todo)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id string) error {
	return s.repo.DeleteTodoById(ctx, id)
}

func (s *TodoService) GetTodosByStatus(ctx context.Context, status string) ([]domain.ToDo, error) {
	return s.repo.GetTodosByStatus(ctx, status)
}
func (s *TodoService) CompleteTodoById(ctx context.Context, id string) error {
	return s.repo.CompleteTodoById(ctx, id)
}
func (s *TodoService) GetTodoByPeriod(ctx context.Context, start string, end string) ([]domain.ToDo, error) {
	return s.repo.GetTodoByPeriod(ctx, start, end)
}
