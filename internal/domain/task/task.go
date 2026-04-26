package task

import (
	"fmt"
	"time"

	"example.com/taskservice/internal/domain"
	"github.com/Krokozabra213/common/types"
)

type Status string

func NewStatus(raw string) (Status, error) {
	status := Status(raw)
	switch status {
	case StatusNew, StatusCompleted, StatusCancelled:
		return status, nil
	default:
		return "", fmt.Errorf("invalid status: %s", raw)
	}
}

const (
	StatusNew       Status = "new"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

type Task struct {
	ID          int64     `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description string    `json:"description" db:"description"`
	Status      Status    `json:"status" db:"status"`
	ScheduledAt time.Time `json:"scheduled_at" db:"scheduled_at"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type CreateInput struct {
	Title       types.NonEmptyString
	Description string
	Status      Status
	ScheduledAt domain.NonEmptyTime
}

func NewCreateInput(title, desc, status, scheduledAt string) (*CreateInput, error) {
	titleVal, err := types.NewNonEmptyString(title)
	if err != nil {
		return nil, fmt.Errorf("invalid title: %w", err)
	}

	statusVal, err := NewStatus(status)
	if err != nil {
		return nil, fmt.Errorf("invalid status: %w", err)
	}

	scheduledVal, err := domain.NewNonEmptyTime(scheduledAt)
	if err != nil {
		return nil, fmt.Errorf("invalid scheduled_at: %w", err)
	}

	return &CreateInput{
		Title:       titleVal,
		Description: desc,
		Status:      statusVal,
		ScheduledAt: scheduledVal,
	}, nil
}

type UpdateInput struct {
	Title       types.NonEmptyString
	Description string
	Status      Status
	ScheduledAt domain.NonEmptyTime
}

func NewUpdateInput(title, desc, status, scheduledAt string) (*UpdateInput, error) {
	titleVal, err := types.NewNonEmptyString(title)
	if err != nil {
		return nil, fmt.Errorf("invalid title: %w", err)
	}

	statusVal, err := NewStatus(status)
	if err != nil {
		return nil, fmt.Errorf("invalid status: %w", err)
	}

	scheduledVal, err := domain.NewNonEmptyTime(scheduledAt)
	if err != nil {
		return nil, fmt.Errorf("invalid scheduled_at: %w", err)
	}

	return &UpdateInput{
		Title:       titleVal,
		Description: desc,
		Status:      statusVal,
		ScheduledAt: scheduledVal,
	}, nil
}
