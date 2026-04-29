package task

import (
	"context"
	"fmt"
	"time"

	"example.com/taskservice/internal/domain"
	taskdomain "example.com/taskservice/internal/domain/task"
	domainTypes "example.com/taskservice/internal/domain/types"
)

type Service struct {
	dbRepository DBRepository
	now          func() time.Time
}

func NewService(repo DBRepository) *Service {
	return &Service{
		dbRepository: repo,
		now:          func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) CreateTask(ctx context.Context, input *taskdomain.CreateInput) (*taskdomain.Task, error) {
	now := s.now()
	model := &taskdomain.Task{
		Title:       input.Title.Value(),
		Description: input.Description,
		Status:      input.Status,
		ScheduledAt: input.ScheduledAt.Value(),
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if input.RecurrenceType == taskdomain.RecurrenceNone {
		return s.dbRepository.CreateTask(ctx, model)
	}

	tx, err := s.dbRepository.BeginTx(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	ruleID, err := s.dbRepository.CreateRecurrenceRuleWithTx(ctx, tx, &taskdomain.Rule{
		Title:          input.Title.Value(),
		Description:    input.Description,
		RecurrenceType: input.RecurrenceType,
		Settings:       input.Settings,
		ScheduledTime:  input.ScheduledAt.Value().Format("15:04:05"),
		Timezone:       input.ScheduledAt.Value().Location().String(),
		CreatedAt:      now,
	})
	if err != nil {
		return nil, err
	}

	model.RecurrenceRuleID = &ruleID
	created, err := s.dbRepository.CreateTaskWithTx(ctx, tx, model)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, domain.NewInternal(domain.MessageInternal, fmt.Errorf("Service.CreateTaskWithRule: failed to commit transaction: %w", err))
	}

	return created, nil
}

func (s *Service) GetTaskByID(ctx context.Context, id domainTypes.PositiveInt[int64]) (*taskdomain.Task, *taskdomain.Rule, error) {
	task, err := s.dbRepository.GetTaskByID(ctx, id.Value())
	if err != nil {
		return nil, nil, err
	}

	var rule *taskdomain.Rule
	if task.RecurrenceRuleID != nil {
		rule, err = s.dbRepository.GetRuleByID(ctx, *task.RecurrenceRuleID)
		if err != nil {
			return nil, nil, err
		}
	}

	return task, rule, nil
}

func (s *Service) UpdateTask(ctx context.Context, id domainTypes.PositiveInt[int64], input *taskdomain.UpdateInput) (*taskdomain.Task, error) {
	model := &taskdomain.Task{
		ID:          id.Value(),
		Title:       input.Title.Value(),
		Description: input.Description,
		Status:      input.Status,
		ScheduledAt: input.ScheduledAt.Value(),
		UpdatedAt:   s.now(),
	}

	updated, err := s.dbRepository.UpdateTask(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) DeleteTask(ctx context.Context, id domainTypes.PositiveInt[int64]) error {
	return s.dbRepository.DeleteTask(ctx, id.Value())
}

func (s *Service) ListTask(ctx context.Context) ([]taskdomain.Task, error) {
	return s.dbRepository.ListTask(ctx)
}

func (s *Service) UpdateRule(ctx context.Context, id domainTypes.PositiveInt[int64], input *taskdomain.UpdateRuleInput) (*taskdomain.Rule, error) {
	model := &taskdomain.Rule{
		ID:             id.Value(),
		Title:          input.Title.Value(),
		Description:    input.Description,
		RecurrenceType: input.RecurrenceType,
		Settings:       input.Settings,
		ScheduledTime:  input.ScheduledAt.Value().Format("15:04:05"),
		Timezone:       input.ScheduledAt.Value().Location().String(),
	}

	updated, err := s.dbRepository.UpdateRule(ctx, model)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

func (s *Service) ListRule(ctx context.Context) ([]taskdomain.Rule, error) {
	return s.dbRepository.ListRule(ctx)
}

func (s *Service) GetRuleByID(ctx context.Context, id domainTypes.PositiveInt[int64]) (*taskdomain.Rule, error) {
	return s.dbRepository.GetRuleByID(ctx, id.Value())
}

func (s *Service) DeleteRule(ctx context.Context, id domainTypes.PositiveInt[int64]) error {
	return s.dbRepository.DeleteRule(ctx, id.Value())
}
