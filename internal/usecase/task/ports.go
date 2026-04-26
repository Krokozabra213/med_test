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

type Usecase interface {
	Create(ctx context.Context, input CreateInput) (*taskdomain.Task, error)
	GetByID(ctx context.Context, id int64) (*taskdomain.Task, error)
	Update(ctx context.Context, id int64, input UpdateInput) (*taskdomain.Task, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]taskdomain.Task, error)
}
