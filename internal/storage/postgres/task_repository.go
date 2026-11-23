package postgres

import (
	"buckwheat-queue/internal/task"
	"context"
	"errors"
	"fmt"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5"
)

type TaskRepository struct {
	db *DB
}

func NewTaskRepository(db *DB) *TaskRepository {
	return &TaskRepository{
		db: db,
	}
}

func (tr *TaskRepository) Create(ctx context.Context, t *task.Task) error {
	sql := `
		INSERT INTO task (
			id, type, payload, status, attempts, max_attempts, timeout, scheduled_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
	`

	_, err := tr.db.Pool.Exec(
		ctx,
		sql,
		t.ID, t.Type, t.Payload, t.Status, t.Attempts, t.MaxAttempts, int(t.Timeout.Seconds()), t.ScheduledAt,
		t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("Insert task: %w", err)
	}

	return nil
}

func (tr *TaskRepository) GetByID(ctx context.Context, id uuid.UUID) (*task.Task, error) {
	sql := `
		SELECT
			id, type, payload, status, attempts, max_attempts, timeout, scheduled_at, created_at, updated_at,
			started_at, completed_at, result, error
		FROM task
		WHERE id = $1
	`

	var t task.Task
	var timeout int
	err := tr.db.Pool.QueryRow(ctx, sql, id).Scan(
		&t.ID, &t.Type, &t.Payload, &t.Status, &t.Attempts, &t.MaxAttempts, &timeout, &t.ScheduledAt, &t.CreatedAt,
		&t.UpdatedAt, &t.StartedAt, &t.CompletedAt, &t.Result, &t.Error,
	)
	t.Timeout = time.Duration(timeout) * time.Second

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Get task: %w", err)
	}

	return &t, nil
}

func (tr *TaskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql := `DELETE FROM task WHERE id = $1`

	_, err := tr.db.Pool.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("Delete task: %w", err)
	}

	return nil
}

func (tr *TaskRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	t, err := tr.GetByID(ctx, id)
	return t != nil, err
}
