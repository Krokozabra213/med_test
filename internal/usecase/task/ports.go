package task

import (
	"context"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	GetTaskByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	UpdateTask(ctx context.Context, task *taskdomain.Task) (*taskdomain.Task, error)
	DeleteTask(ctx context.Context, id int64) error
	ListTask(ctx context.Context) ([]taskdomain.Task, error)
}
