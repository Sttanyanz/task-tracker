package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var ErrTaskNotFound = errors.New("task not found")

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(ctx context.Context, db *sql.DB) (*PostgresStorage, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}

	const createTasksTable = `
		CREATE TABLE IF NOT EXISTS tasks (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			name TEXT NOT NULL,
			description TEXT NOT NULL DEFAULT '',
			completed BOOLEAN NOT NULL DEFAULT FALSE
		)`
	if _, err := db.ExecContext(ctx, createTasksTable); err != nil {
		return nil, fmt.Errorf("create tasks table: %w", err)
	}

	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) CreateTask(ctx context.Context, task Task) (int64, error) {
	const query = `
		INSERT INTO tasks (user_id, name, description, completed)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	var id int64
	err := s.db.QueryRowContext(ctx, query, task.UserID, task.Name, task.Description, task.Completed).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("create task: %w", err)
	}
	return id, nil
}

func (s *PostgresStorage) GetTask(ctx context.Context, id int64) (Task, error) {
	const query = `
		SELECT id, user_id, name, description, completed
		FROM tasks
		WHERE id = $1`

	var task Task
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.UserID,
		&task.Name,
		&task.Description,
		&task.Completed,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, fmt.Errorf("%w: %d", ErrTaskNotFound, id)
	}
	if err != nil {
		return Task{}, fmt.Errorf("get task %d: %w", id, err)
	}
	return task, nil
}

func (s *PostgresStorage) UpdateTask(ctx context.Context, id int64, task Task) error {
	const query = `
		UPDATE tasks
		SET user_id = $1, name = $2, description = $3, completed = $4
		WHERE id = $5`

	result, err := s.db.ExecContext(ctx, query, task.UserID, task.Name, task.Description, task.Completed, id)
	if err != nil {
		return fmt.Errorf("update task %d: %w", id, err)
	}
	if err := requireTaskAffected(result, id); err != nil {
		return err
	}
	return nil
}

func (s *PostgresStorage) DeleteTask(ctx context.Context, id int64) error {
	result, err := s.db.ExecContext(ctx, "DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete task %d: %w", id, err)
	}
	if err := requireTaskAffected(result, id); err != nil {
		return err
	}
	return nil
}

func requireTaskAffected(result sql.Result, id int64) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check affected rows for task %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w: %d", ErrTaskNotFound, id)
	}
	return nil
}
