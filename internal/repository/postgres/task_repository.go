package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	taskdomain "example.com/taskservice/internal/domain/task"
)

const (
	ctxTimeout = 3 * time.Second
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateTask(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	query := `
			INSERT INTO tasks (title, description, status, scheduled_at, created_at, updated_at)
			VALUES (@title, @description, @status, @scheduled_at, @created_at, @updated_at)
			RETURNING id, title, description, status, scheduled_at, created_at, updated_at
	`

	args := pgx.NamedArgs{
		"title":        task.Title,
		"description":  task.Description,
		"status":       task.Status,
		"scheduled_at": task.ScheduledAt,
		"created_at":   task.CreatedAt,
		"updated_at":   task.UpdatedAt,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("repository.CreateTask: query failed: %w", err)
	}
	defer rows.Close()

	created, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[taskdomain.Task])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repository.CreateTask: no rows returned: %w", err)
		}
		return nil, fmt.Errorf("repository.CreateTask: collect row failed: %w", err)
	}

	return &created, nil
}

func (r *Repository) GetTaskByID(ctx context.Context, id int64) (*taskdomain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	query := `
        SELECT id, title, description, status, scheduled_at, created_at, updated_at
        FROM tasks
        WHERE id = @id
    `

	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("repository.GetTaskByID: query failed: %w", err)
	}
	defer rows.Close()

	task, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[taskdomain.Task])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("repository.GetTaskByID: collect row failed: %w", err)
	}

	return &task, nil
}

func (r *Repository) UpdateTask(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	query := `
        UPDATE tasks
        SET title = @title,
            description = @description,
            status = @status,
            scheduled_at = @scheduled_at,
            updated_at = @updated_at
        WHERE id = @id
        RETURNING id, title, description, status, scheduled_at, created_at, updated_at
    `

	args := pgx.NamedArgs{
		"id":           task.ID,
		"title":        task.Title,
		"description":  task.Description,
		"status":       task.Status,
		"scheduled_at": task.ScheduledAt,
		"updated_at":   task.UpdatedAt,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, fmt.Errorf("repository.UpdateTask: query failed: %w", err)
	}
	defer rows.Close()

	updated, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[taskdomain.Task])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, taskdomain.ErrTaskNotFound
		}
		return nil, fmt.Errorf("repository.UpdateTask: collect row failed: %w", err)
	}

	return &updated, nil
}

func (r *Repository) DeleteTask(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	query := `DELETE FROM tasks WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	result, err := r.pool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("repository.DeleteTask: execute failed: %w", err)
	}

	if result.RowsAffected() == 0 {
		return taskdomain.ErrTaskNotFound
	}

	return nil
}

func (r *Repository) ListTask(ctx context.Context) ([]taskdomain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	query := `
        SELECT id, title, description, status, scheduled_at, created_at, updated_at
        FROM tasks
        ORDER BY id DESC
    `

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("repository.ListTask: query failed: %w", err)
	}
	defer rows.Close()

	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[taskdomain.Task])
	if err != nil {
		return nil, fmt.Errorf("repository.ListTask: collect rows failed: %w", err)
	}

	return tasks, nil
}
