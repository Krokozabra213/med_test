package task

import (
	"context"

	taskdomain "example.com/taskservice/internal/domain/task"
	"github.com/jackc/pgx/v5"
)

type DBRepository interface {
	CreateTask(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	UpdateTask(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	DeleteTask(ctx context.Context, id int64) error
	ListTask(ctx context.Context) ([]taskdomain.Task, error)

	GetRuleByID(ctx context.Context, id int64) (*taskdomain.Rule, error)
	UpdateRule(ctx context.Context, rule *taskdomain.Rule) (*taskdomain.Rule, error)
	DeleteRule(ctx context.Context, id int64) error
	ListRule(ctx context.Context) ([]taskdomain.Rule, error)

	BeginTx(ctx context.Context) (pgx.Tx, error)
	CreateTaskWithTx(ctx context.Context, tx pgx.Tx, task *taskdomain.Task) (*taskdomain.Task, error)
	CreateRecurrenceRuleWithTx(ctx context.Context, tx pgx.Tx, rule *taskdomain.Rule) (int64, error)
}
