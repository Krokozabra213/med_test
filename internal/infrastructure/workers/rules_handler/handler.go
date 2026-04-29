package ruleshandler

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type repository interface {
	ListRule(ctx context.Context) ([]taskdomain.Rule, error)
	CreateGeneratedTask(ctx context.Context, task *taskdomain.Task, ruleID int64) error
}

type RecurrenceHandler interface {
	CanHandle(rt taskdomain.RecurrenceType) bool
	Handle(rule *taskdomain.Rule, date time.Time) (bool, error)
}

type RecurrenceWorker struct {
	log      *slog.Logger
	repo     repository
	handlers []RecurrenceHandler
	now      func() time.Time
}

func NewRecurrenceWorker(
	log *slog.Logger,
	repo repository,
	handlers ...RecurrenceHandler,
) *RecurrenceWorker {
	return &RecurrenceWorker{
		log:      log,
		repo:     repo,
		handlers: handlers,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

func (w *RecurrenceWorker) RunOnce(ctx context.Context) error {

	rules, err := w.repo.ListRule(ctx)
	if err != nil {
		return err
	}

	for _, rule := range rules {
		if err := w.processRule(ctx, &rule); err != nil {
			w.log.Warn("RecurrenceWorker.processRule", "error", err.Error())
		}
	}

	return nil
}

func (w *RecurrenceWorker) Run(ctx context.Context, period time.Duration) {
	ticker := time.NewTicker(period)
	defer ticker.Stop()

	var err error

	for {
		// для приоритизации
		select {
		case <-ctx.Done():
			return
		default:
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err = w.RunOnce(ctx); err != nil {
				w.log.Warn("RecurrenceWorker.RunOnce", "error", err.Error())
			}
		}
	}
}

func (w *RecurrenceWorker) processRule(ctx context.Context, rule *taskdomain.Rule) error {

	loc, err := time.LoadLocation(rule.Timezone)
	if err != nil {
		return err
	}

	now := w.now().In(loc)

	parsedTime, err := time.Parse("15:04:05", rule.ScheduledTime)
	if err != nil {
		return err
	}

	scheduledAt := time.Date(
		now.Year(),
		now.Month(),
		now.Day(),
		parsedTime.Hour(),
		parsedTime.Minute(),
		parsedTime.Second(),
		0,
		loc,
	)

	if now.Before(scheduledAt) {
		return nil
	}

	handler := w.getHandler(rule.RecurrenceType)
	if handler == nil {
		return fmt.Errorf("no handler for recurrence type %s", rule.RecurrenceType)
	}

	match, err := handler.Handle(rule, scheduledAt)
	if err != nil {
		return err
	}

	if !match {
		return nil
	}

	taskNow := w.now()

	task := &taskdomain.Task{
		Title:       rule.Title,
		Description: rule.Description,
		Status:      taskdomain.StatusNew,
		ScheduledAt: scheduledAt,
		CreatedAt:   taskNow,
		UpdatedAt:   taskNow,
	}

	return w.repo.CreateGeneratedTask(ctx, task, rule.ID)
}

func (w *RecurrenceWorker) getHandler(rt taskdomain.RecurrenceType) RecurrenceHandler {
	for _, h := range w.handlers {
		if h.CanHandle(rt) {
			return h
		}
	}
	return nil
}
