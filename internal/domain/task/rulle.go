package task

import (
	"fmt"
	"time"

	domainTypes "example.com/taskservice/internal/domain/types"
)

type Rule struct {
	ID             int64          `db:"id"`
	Title          string         `db:"title"`
	Description    string         `db:"description"`
	RecurrenceType RecurrenceType `db:"recurrence_type"`
	Settings       []byte         `db:"settings"`
	ScheduledTime  string         `db:"scheduled_time"`
	Timezone       string         `db:"timezone"`
	CreatedAt      time.Time      `db:"created_at"`
}

type UpdateRuleInput struct {
	Title          domainTypes.NonEmptyString
	Description    string
	RecurrenceType RecurrenceType
	Settings       []byte
	ScheduledAt    domainTypes.NonEmptyTime
}

func NewUpdateRuleInput(title, desc, scheduledAt, recurrenceType string, settings []byte) (*UpdateRuleInput, error) {
	titleVal, err := domainTypes.NewNonEmptyString(title)
	if err != nil {
		return nil, fmt.Errorf("invalid title: %w", err)
	}

	scheduledVal, err := domainTypes.NewNonEmptyTime(scheduledAt)
	if err != nil {
		return nil, fmt.Errorf("invalid scheduled_at: %w", err)
	}

	recurrenceVal, err := NewRecurrenceType(recurrenceType)
	if err != nil {
		return nil, err
	}

	if err = validateSettings(recurrenceVal, settings); err != nil {
		return nil, err
	}

	return &UpdateRuleInput{
		Title:          titleVal,
		Description:    desc,
		RecurrenceType: recurrenceVal,
		Settings:       settings,
		ScheduledAt:    scheduledVal,
	}, nil
}
