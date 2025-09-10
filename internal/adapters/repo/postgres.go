package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ToDo-List/internal/core/domain"
	"ToDo-List/internal/core/ports"
)

type PostgreRepo struct {
	db *sql.DB
}

func NewPostgreRepo(db *sql.DB) ports.PostgreRepo {
	return &PostgreRepo{db: db}
}

func (r *PostgreRepo) GetAllTodos(ctx context.Context) ([]domain.ToDo, error) {
	query := `SELECT id, todo, message, created_at, updated_at, deadline, system_message, completed_at, complete FROM todo`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []domain.ToDo

	for rows.Next() {
		var todo domain.ToDo
		err := rows.Scan(
			&todo.Id,
			&todo.Todo,
			&todo.Message,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&todo.Deadline,
			&todo.SystemMessage,
			&todo.CompletedAt,
			&todo.Complete,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (r *PostgreRepo) GetTodoById(ctx context.Context, id string) (domain.ToDo, error) {
	query := `SELECT id, todo, message, created_at, updated_at, deadline, system_message, completed_at, complete 
	          FROM todo WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var todo domain.ToDo
	err := row.Scan(
		&todo.Id,
		&todo.Todo,
		&todo.Message,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.Deadline,
		&todo.SystemMessage,
		&todo.CompletedAt,
		&todo.Complete,
	)
	if err != nil {
		return domain.ToDo{}, err
	}

	return todo, nil
}

func (r *PostgreRepo) DeleteTodoById(ctx context.Context, id string) error {
	query := `DELETE FROM todo WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %s not found", id)
	}

	return nil
}

func (r *PostgreRepo) UpdateTodo(ctx context.Context, todo domain.ToDo) error {
	query := `
		UPDATE todo 
		SET 
			todo = $1,
			message = $2,
			updated_at = $3,
			deadline = $4,
			system_message = $5,
			completed_at = $6,
			complete = $7
		WHERE id = $8
	`

	todo.UpdatedAt = time.Now()

	result, err := r.db.ExecContext(ctx, query,
		todo.Todo,
		todo.Message,
		todo.UpdatedAt,
		todo.Deadline,
		todo.SystemMessage,
		todo.CompletedAt,
		todo.Complete,
		todo.Id,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo with id %s not found", todo.Id)
	}

	return nil
}

func (r *PostgreRepo) CreateTodo(ctx context.Context, todo domain.ToDo) (domain.ToDo, error) {
	query := `
		INSERT INTO todo (
			id, todo, message, created_at, updated_at, deadline, system_message, completed_at, complete
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.ExecContext(ctx, query,
		todo.Id,
		todo.Todo,
		todo.Message,
		todo.CreatedAt,
		todo.UpdatedAt,
		todo.Deadline,
		todo.SystemMessage,
		todo.CompletedAt,
		todo.Complete,
	)
	if err != nil {
		return domain.ToDo{}, err
	}

	return todo, nil
}
