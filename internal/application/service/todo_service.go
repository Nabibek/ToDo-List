package service

import (
	"context"
	"time"

	"ToDo-List/internal/adapters/logger"
	"ToDo-List/internal/core/domain"
	"ToDo-List/internal/core/ports"
)

type TodoService struct {
	repo   ports.PostgreRepo
	logger *logger.Logger
}

func NewToDoService(repo ports.PostgreRepo, logger *logger.Logger) ports.ToDoService {
	return &TodoService{
		repo:   repo,
		logger: logger,
	}
}

func (s *TodoService) CreateTodo(ctx context.Context, todo domain.ToDo) (domain.ToDo, error) {
	s.logger.Debug("Creating todo in service")
	return s.repo.CreateTodo(ctx, todo)
}
func (s *TodoService) GetAllTodosWithFilters(ctx context.Context, filter ports.TodoFilter) ([]domain.ToDo, error) {
	s.logger.Debug("Getting all todos with filters: %+v", filter)
	return s.repo.GetAllTodosWithFilters(ctx, filter)
}

func (s *TodoService) GetTodoById(ctx context.Context, id string) (domain.ToDo, error) {
	s.logger.Debug("Getting todo by ID: %s", id)
	return s.repo.GetTodoById(ctx, id)
}

func (s *TodoService) UpdateTodo(ctx context.Context, todo domain.ToDo) error {
	s.logger.Debug("Updating todo: %s", todo.Id)
	todo.UpdatedAt = time.Now()
	return s.repo.UpdateTodo(ctx, todo)
}

func (s *TodoService) DeleteTodo(ctx context.Context, id string) error {
	s.logger.Debug("Deleting todo: %s", id)
	return s.repo.DeleteTodoById(ctx, id)
}

func (s *TodoService) CompleteTodoById(ctx context.Context, id string) error {
	s.logger.Debug("Completing todo: %s", id)
	
	
	
	todo, err := s.repo.GetTodoById(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get todo for completion: %s, error: %v", id, err)
		return err
	}
	if todo.Complete {
		s.logger.Warn("Todo %s is already completed", id)
		return nil 
	}

	todo.Complete = true
	todo.CompletedAt = time.Now()
	todo.UpdatedAt = time.Now()

	s.logger.Info("Marking todo as completed: %s", id)
	return s.repo.UpdateTodo(ctx, todo)
}