package postgres

import (
	"context"
	"errors"
	"fmt"

	"example.com/taskservice/internal/domain"
	taskdomain "example.com/taskservice/internal/domain/task"
	"github.com/jackc/pgx/v5"
)

const ruleEntity = "recurrence rule"

func (r *Repository) CreateRecurrenceRuleWithTx(ctx context.Context, tx pgx.Tx, rule *taskdomain.Rule) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	query := `
			INSERT INTO task_recurrence_rules
					(title, description, recurrence_type, settings, scheduled_time, timezone, created_at)
			VALUES
					(@title, @description, @recurrence_type, @settings, @scheduled_time, @timezone, @created_at)
			RETURNING id
	`

	args := pgx.NamedArgs{
		"title":           rule.Title,
		"description":     rule.Description,
		"recurrence_type": rule.RecurrenceType,
		"settings":        rule.Settings,
		"scheduled_time":  rule.ScheduledTime,
		"timezone":        rule.Timezone,
		"created_at":      rule.CreatedAt,
	}

	var ruleID int64
	err := tx.QueryRow(ctx, query, args).Scan(&ruleID)
	if err != nil {
		return 0, domain.NewInternal(domain.MessageInternal, fmt.Errorf("repository.CreateRecurrenceRuleWithTx: query failed: %w", err))
	}
	return ruleID, nil
}

func (r *Repository) GetRuleByID(ctx context.Context, id int64) (*taskdomain.Rule, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	query := `
		SELECT
			id,
			title,
			description,
			recurrence_type,
			settings,
			scheduled_time,
			timezone,
			created_at
		FROM task_recurrence_rules
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id": id,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, domain.NewInternal(domain.MessageInternal, fmt.Errorf("repository.GetRuleByID: query failed: %w", err))
	}
	defer rows.Close()

	rule, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[taskdomain.Rule])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewNotFound(ruleEntity)
		}
		return nil, domain.NewInternal(domain.MessageInternal, fmt.Errorf("repository.GetRuleByID: collect failed: %w", err))
	}

	return &rule, nil
}

func (r *Repository) UpdateRule(ctx context.Context, rule *taskdomain.Rule) (*taskdomain.Rule, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	query := `
			UPDATE task_recurrence_rules
			SET title = @title,
					description = @description,
					recurrence_type = @recurrence_type,
					settings = @settings,
					scheduled_time = @scheduled_time,
					timezone = @timezone
			WHERE id = @id
			RETURNING id, title, description, recurrence_type, settings, scheduled_time, timezone, created_at
	`

	args := pgx.NamedArgs{
		"id":              rule.ID,
		"title":           rule.Title,
		"description":     rule.Description,
		"recurrence_type": string(rule.RecurrenceType),
		"settings":        rule.Settings,
		"scheduled_time":  rule.ScheduledTime,
		"timezone":        rule.Timezone,
	}

	rows, err := r.pool.Query(ctx, query, args)
	if err != nil {
		return nil, domain.NewInternal(domain.MessageInternal, fmt.Errorf("repository.UpdateRule: query failed: %w", err))
	}
	defer rows.Close()

	updated, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[taskdomain.Rule])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.NewInternal(domain.MessageInternal, fmt.Errorf("repository.UpdateRule: collect row not found: %w", err))
		}
		return nil, domain.NewInternal(domain.MessageInternal, fmt.Errorf("repository.UpdateRule: collect row failed: %w", err))
	}

	return &updated, nil
}

func (r *Repository) DeleteRule(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	query := `DELETE FROM task_recurrence_rules WHERE id = @id`

	args := pgx.NamedArgs{
		"id": id,
	}

	result, err := r.pool.Exec(ctx, query, args)
	if err != nil {
		return domain.NewInternal(domain.MessageInternal, fmt.Errorf("repository.DeleteRule: execute failed: %w", err))
	}

	if result.RowsAffected() == 0 {
		return domain.NewNotFound(ruleEntity)
	}

	return nil
}

func (r *Repository) ListRule(ctx context.Context) ([]taskdomain.Rule, error) {
	ctx, cancel := context.WithTimeout(ctx, ctxTimeout)
	defer cancel()

	query := `
        SELECT id, title, description, recurrence_type, settings, scheduled_time, timezone, created_at
        FROM task_recurrence_rules
        ORDER BY id DESC
    `

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, domain.NewInternal(domain.MessageInternal, fmt.Errorf("repository.ListRule: query failed: %w", err))
	}
	defer rows.Close()

	rules, err := pgx.CollectRows(rows, pgx.RowToStructByName[taskdomain.Rule])
	if err != nil {
		return nil, domain.NewInternal(domain.MessageInternal, fmt.Errorf("repository.ListRule: collect rows failed: %w", err))
	}

	return rules, nil
}
