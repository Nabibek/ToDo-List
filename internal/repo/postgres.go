package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"
"strings"
	"ToDo-List/internal/adapters/logger"
	"ToDo-List/internal/core/domain"
	"ToDo-List/internal/core/ports"
)

type PostgreRepo struct {
	db     *sql.DB
	logger *logger.Logger
}

func NewPostgreRepo(db *sql.DB, logger *logger.Logger) ports.PostgreRepo {
	return &PostgreRepo{
		db:     db,
		logger: logger,
	}
}

func (r *PostgreRepo) GetAllTodosWithFilters(ctx context.Context, filter ports.TodoFilter) ([]domain.ToDo, error) {
	r.logger.Debug("Executing GetAllTodosWithFilters: %+v", filter)
	
	query := `SELECT id, todo, message, created_at, updated_at, deadline, priority, completed_at, complete 
              FROM todo`
	
	args := []interface{}{}
	conditions := []string{}
	
	// Фильтрация по статусу
	switch filter.Status {
	case "active":
		conditions = append(conditions, "complete = false")
	case "completed":
		conditions = append(conditions, "complete = true")
	case "overdue":
		now := time.Now().Format("2006-01-02 15:04:05")
		conditions = append(conditions, "complete = false AND deadline < $1")
		args = append(args, now)
	}
	
	// Фильтрация по периоду
	switch filter.Period {
	case "today":
		today := time.Now().Format("2006-01-02")
		conditions = append(conditions, "DATE(created_at) = $"+fmt.Sprint(len(args)+1))
		args = append(args, today)
	case "week":
		weekAgo := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
		conditions = append(conditions, "created_at >= $"+fmt.Sprint(len(args)+1))
		args = append(args, weekAgo)
	case "month":
		monthAgo := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
		conditions = append(conditions, "created_at >= $"+fmt.Sprint(len(args)+1))
		args = append(args, monthAgo)
	}
	
	// Добавляем условия WHERE
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	
	// Сортировка
	orderBy := "created_at"
	if filter.OrderBy != "" {
		orderBy = filter.OrderBy
	}
	
	orderDir := "DESC"
	if filter.OrderDir == "asc" {
		orderDir = "ASC"
	}
	
	if orderBy == "priority" {
		query += ` ORDER BY 
			CASE priority 
				WHEN 'high' THEN 1 
				WHEN 'medium' THEN 2 
				WHEN 'low' THEN 3 
				ELSE 4 
			END ` + orderDir
	} else {
		query += " ORDER BY " + orderBy + " " + orderDir
	}
	
	r.logger.Debug("SQL Query: %s, Args: %v", query, args)
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Query failed: %v", err)
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
			&todo.Priority,
			&todo.CompletedAt,
			&todo.Complete,
		)
		if err != nil {
			r.logger.Error("Row scan failed: %v", err)
			return nil, err
		}
		todos = append(todos, todo)
	}
	
	if err := rows.Err(); err != nil {
		r.logger.Error("Rows error: %v", err)
		return nil, err
	}
	
	r.logger.Info("Retrieved %d todos", len(todos))
	return todos, nil
}

func (r *PostgreRepo) GetTodoById(ctx context.Context, id string) (domain.ToDo, error) {
	r.logger.Debug("Executing GetTodoById: id=%s", id)
	
	query := `SELECT id, todo, message, created_at, updated_at, deadline, priority, completed_at, complete 
	          FROM todo WHERE id = $1`

	r.logger.Debug("SQL Query: %s, Arg: %s", query, id)
	row := r.db.QueryRowContext(ctx, query, id)

	var todo domain.ToDo
	err := row.Scan(
		&todo.Id,
		&todo.Todo,
		&todo.Message,
		&todo.CreatedAt,
		&todo.UpdatedAt,
		&todo.Deadline,
		&todo.Priority,
		&todo.CompletedAt,
		&todo.Complete,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			r.logger.Warn("Todo not found: %s", id)
			return domain.ToDo{}, fmt.Errorf("todo not found: %s", id)
		}
		r.logger.Error("Scan failed: %v", err)
		return domain.ToDo{}, err
	}

	r.logger.Debug("Todo found: %s", id)
	return todo, nil
}

func (r *PostgreRepo) DeleteTodoById(ctx context.Context, id string) error {
	r.logger.Debug("Executing DeleteTodoById: id=%s", id)
	
	query := `DELETE FROM todo WHERE id = $1`
	r.logger.Debug("SQL Query: %s, Arg: %s", query, id)

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Delete failed: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("RowsAffected failed: %v", err)
		return err
	}

	if rowsAffected == 0 {
		r.logger.Warn("Todo not found for deletion: %s", id)
		return fmt.Errorf("todo with id %s not found", id)
	}

	r.logger.Info("Todo deleted successfully: %s", id)
	return nil
}

func (r *PostgreRepo) UpdateTodo(ctx context.Context, todo domain.ToDo) error {
	r.logger.Debug("Executing UpdateTodo: id=%s", todo.Id)
	
	query := `
		UPDATE todo 
		SET 
			todo = $1,
			message = $2,
			updated_at = $3,
			deadline = $4,
			priority = $5,
			completed_at = $6,
			complete = $7
		WHERE id = $8
	`

	todo.UpdatedAt = time.Now()

	r.logger.Debug("SQL Query: %s, Args: %+v", query, todo)
	result, err := r.db.ExecContext(ctx, query,
		todo.Todo,
		todo.Message,
		todo.UpdatedAt,
		todo.Deadline,
		todo.Priority,
		todo.CompletedAt,
		todo.Complete,
		todo.Id,
	)
	if err != nil {
		r.logger.Error("Update failed: %v", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("RowsAffected failed: %v", err)
		return err
	}

	if rowsAffected == 0 {
		r.logger.Warn("Todo not found for update: %s", todo.Id)
		return fmt.Errorf("todo with id %s not found", todo.Id)
	}

	r.logger.Info("Todo updated successfully: %s", todo.Id)
	return nil
}

func (r *PostgreRepo) CreateTodo(ctx context.Context, todo domain.ToDo) (domain.ToDo, error) {
	r.logger.Debug("Executing CreateTodo: %+v", todo)
	
	query := `
		INSERT INTO todo (
			id, todo, message, created_at, updated_at, deadline, priority, completed_at, complete
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	r.logger.Debug("SQL Query: %s, Args: %+v", query, todo)
	_, err := r.db.ExecContext(ctx, query,
		todo.Id,
		todo.Todo,
		todo.Message,
		todo.CreatedAt,
		todo.UpdatedAt,
		todo.Deadline,
		todo.Priority,
		todo.CompletedAt,
		todo.Complete,
	)
	if err != nil {
		r.logger.Error("Insert failed: %v", err)
		return domain.ToDo{}, err
	}

	r.logger.Info("Todo created successfully: %s", todo.Id)
	return todo, nil
}



func (r *PostgreRepo) Ping() error {
	r.logger.Debug("Pinging database")
	err := r.db.Ping()
	if err != nil {
		r.logger.Error("Database ping failed: %v", err)
		return err
	}
	r.logger.Debug("Database ping successful")
	return nil
}