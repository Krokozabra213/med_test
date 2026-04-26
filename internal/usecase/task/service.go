package task

import (
	"context"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
	"github.com/Krokozabra213/common/types"
)

type Service struct {
	taskRepo TaskRepository
	now      func() time.Time
}

func NewService(repo TaskRepository) *Service {
	return &Service{
		taskRepo: repo,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, input *taskdomain.CreateInput) (*taskdomain.Task, error) {

	model := &taskdomain.Task{
		Title:       input.Title.Value(),
		Description: input.Description,
		Status:      input.Status,
		ScheduledAt: input.ScheduledAt.Value(),
	}
	now := s.now()
	model.CreatedAt = now
	model.UpdatedAt = now

	created, err := s.taskRepo.CreateTask(ctx, model)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetByID(ctx context.Context, id types.PositiveInt[int64]) (*taskdomain.Task, error) {
	return s.taskRepo.GetTaskByID(ctx, id.Value())
}

func (s *Service) Update(ctx context.Context, id types.PositiveInt[int64], input *taskdomain.UpdateInput) (*taskdomain.Task, error) {
	model := &taskdomain.Task{
		ID:          id.Value(),
		Title:       input.Title.Value(),
		Description: input.Description,
		Status:      input.Status,
		ScheduledAt: input.ScheduledAt.Value(),
		UpdatedAt:   s.now(),
	}

	updated, err := s.taskRepo.UpdateTask(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id types.PositiveInt[int64]) error {
	return s.taskRepo.DeleteTask(ctx, id.Value())
}

func (s *Service) List(ctx context.Context) ([]taskdomain.Task, error) {
	return s.taskRepo.ListTask(ctx)
}
