package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	domainTypes "example.com/taskservice/internal/domain/types"
)

type Status string

const (
	StatusNew       Status = "new"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

func NewStatus(raw string) (Status, error) {
	status := Status(raw)
	switch status {
	case StatusNew, StatusCompleted, StatusCancelled:
		return status, nil
	default:
		return "", fmt.Errorf("invalid status: %s", raw)
	}
}

type RecurrenceType string

const (
	RecurrenceNone       RecurrenceType = "empty"
	RecurrenceDaily      RecurrenceType = "daily"
	RecurrenceMonthlyDay RecurrenceType = "monthly_day"
	RecurrenceSpecific   RecurrenceType = "specific_dates"
	RecurrenceDayParity  RecurrenceType = "day_parity"
)

func NewRecurrenceType(raw string) (RecurrenceType, error) {
	rt := RecurrenceType(raw)

	switch rt {
	case RecurrenceNone,
		RecurrenceDaily,
		RecurrenceMonthlyDay,
		RecurrenceSpecific,
		RecurrenceDayParity:
		return rt, nil
	default:
		return "", fmt.Errorf("invalid recurrence_type: %s", raw)
	}
}

type Task struct {
	ID               int64     `json:"id" db:"id"`
	Title            string    `json:"title" db:"title"`
	Description      string    `json:"description" db:"description"`
	Status           Status    `json:"status" db:"status"`
	ScheduledAt      time.Time `json:"scheduled_at" db:"scheduled_at"`
	RecurrenceRuleID *int64    `json:"recurrence_rule_id" db:"recurrence_rule_id"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

type CreateInput struct {
	Title          domainTypes.NonEmptyString
	Description    string
	Status         Status
	ScheduledAt    domainTypes.NonEmptyTime
	RecurrenceType RecurrenceType
	Settings       []byte
}

func NewCreateInput(title, desc, status, scheduledAt, recurrenceType string, settings []byte) (*CreateInput, error) {
	titleVal, err := domainTypes.NewNonEmptyString(title)
	if err != nil {
		return nil, fmt.Errorf("invalid title: %w", err)
	}

	statusVal, err := NewStatus(status)
	if err != nil {
		return nil, err
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

	return &CreateInput{
		Title:          titleVal,
		Description:    desc,
		Status:         statusVal,
		ScheduledAt:    scheduledVal,
		RecurrenceType: recurrenceVal,
		Settings:       settings,
	}, nil
}

type DailySettings struct {
	Interval int `json:"interval"`
}

type MonthlyDaySettings struct {
	Day int `json:"day"`
}

type SpecificDatesSettings struct {
	Dates []string `json:"dates"`
}

type DayParitySettings struct {
	Parity string `json:"parity"`
}

func validateSettings(rt RecurrenceType, raw []byte) error {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()

	if rt == RecurrenceNone {
		if len(raw) != 0 && string(raw) != "{}" {
			return fmt.Errorf("settings must be empty when recurrence_type is empty")
		}
		return nil
	}

	if len(raw) == 0 {
		return fmt.Errorf("settings required for recurrence_type %s", rt)
	}

	switch rt {

	case RecurrenceDaily:
		var s DailySettings
		if err := dec.Decode(&s); err != nil {
			return fmt.Errorf("invalid daily settings: %w", err)
		}
		if s.Interval <= 0 {
			return fmt.Errorf("interval must be > 0")
		}

	case RecurrenceMonthlyDay:
		var s MonthlyDaySettings
		if err := dec.Decode(&s); err != nil {
			return fmt.Errorf("invalid monthly_day settings: %w", err)
		}
		if s.Day < 1 || s.Day > 30 {
			return fmt.Errorf("day must be between 1 and 30")
		}

	case RecurrenceSpecific:
		var s SpecificDatesSettings
		if err := dec.Decode(&s); err != nil {
			return fmt.Errorf("invalid specific_dates settings: %w", err)
		}
		if len(s.Dates) == 0 {
			return fmt.Errorf("dates cannot be empty")
		}

		for _, d := range s.Dates {
			if _, err := time.Parse("2006-01-02", d); err != nil {
				return fmt.Errorf("invalid date format: %s", d)
			}
		}

	case RecurrenceDayParity:
		var s DayParitySettings
		if err := dec.Decode(&s); err != nil {
			return fmt.Errorf("invalid day_parity settings: %w", err)
		}

		if s.Parity != "even" && s.Parity != "odd" {
			return fmt.Errorf("parity must be 'even' or 'odd'")
		}

	default:
		return fmt.Errorf("unsupported recurrence_type")
	}

	return nil
}

type UpdateInput struct {
	Title       domainTypes.NonEmptyString
	Description string
	Status      Status
	ScheduledAt domainTypes.NonEmptyTime
}

func NewUpdateInput(title, desc, status, scheduledAt string) (*UpdateInput, error) {
	titleVal, err := domainTypes.NewNonEmptyString(title)
	if err != nil {
		return nil, fmt.Errorf("invalid title: %w", err)
	}

	statusVal, err := NewStatus(status)
	if err != nil {
		return nil, fmt.Errorf("invalid status: %w", err)
	}

	scheduledVal, err := domainTypes.NewNonEmptyTime(scheduledAt)
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
