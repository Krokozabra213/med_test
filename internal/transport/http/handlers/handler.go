package handlers

import (
	"context"
	"log/slog"

	taskdomain "example.com/taskservice/internal/domain/task"
	domainTypes "example.com/taskservice/internal/domain/types"
)

type Usecase interface {
	CreateTask(ctx context.Context, input *taskdomain.CreateInput) (*taskdomain.Task, error)
	GetTaskByID(ctx context.Context, id domainTypes.PositiveInt[int64]) (*taskdomain.Task, *taskdomain.Rule, error)
	UpdateTask(ctx context.Context, id domainTypes.PositiveInt[int64], input *taskdomain.UpdateInput) (*taskdomain.Task, error)
	DeleteTask(ctx context.Context, id domainTypes.PositiveInt[int64]) error
	ListTask(ctx context.Context) ([]taskdomain.Task, error)

	UpdateRule(ctx context.Context, id domainTypes.PositiveInt[int64], input *taskdomain.UpdateRuleInput) (*taskdomain.Rule, error)
	ListRule(ctx context.Context) ([]taskdomain.Rule, error)
	GetRuleByID(ctx context.Context, id domainTypes.PositiveInt[int64]) (*taskdomain.Rule, error)
	DeleteRule(ctx context.Context, id domainTypes.PositiveInt[int64]) error
}

type Handler struct {
	log     *slog.Logger
	usecase Usecase
}

func NewHandler(usecase Usecase, log *slog.Logger) *Handler {
	return &Handler{
		usecase: usecase,
		log:     log,
	}
}
