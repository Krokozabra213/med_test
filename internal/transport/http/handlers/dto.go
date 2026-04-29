package handlers

import (
	"encoding/json"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type taskUpdateDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	ScheduledAt string `json:"scheduled_at"`
}

type ruleMutationDTO struct {
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	ScheduledAt    string          `json:"scheduled_at"`
	RecurrenceType string          `json:"recurrence_type"`
	Settings       json.RawMessage `json:"settings"`
}

type taskMutationDTO struct {
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	Status         string          `json:"status"`
	ScheduledAt    string          `json:"scheduled_at"`
	RecurrenceType string          `json:"recurrence_type"`
	Settings       json.RawMessage `json:"settings"`
}

type taskDTO struct {
	ID               int64             `json:"id"`
	Title            string            `json:"title"`
	Description      string            `json:"description"`
	Status           taskdomain.Status `json:"status"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`
	RecurrenceRuleID *int64            `json:"recurrence_rule_id,omitempty"`
	ScheduledAt      string            `json:"scheduled_at"`
}

func newTaskDTO(task *taskdomain.Task) taskDTO {
	if task == nil {
		return taskDTO{}
	}
	return taskDTO{
		ID:               task.ID,
		Title:            task.Title,
		Description:      task.Description,
		Status:           task.Status,
		CreatedAt:        task.CreatedAt,
		UpdatedAt:        task.UpdatedAt,
		RecurrenceRuleID: task.RecurrenceRuleID,
		ScheduledAt:      task.ScheduledAt.Format(time.RFC3339),
	}
}

type taskWithRuleDTO struct {
	Task taskDTO  `json:"task"`
	Rule *ruleDTO `json:"rule,omitempty"`
}

type ruleDTO struct {
	ID             int64           `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	RecurrenceType string          `json:"recurrence_type"`
	Settings       json.RawMessage `json:"settings"`
	ScheduledTime  string          `json:"scheduled_time"`
	Timezone       string          `json:"timezone"`
}

func newRuleDTO(rule *taskdomain.Rule) *ruleDTO {
	if rule == nil {
		return nil
	}

	return &ruleDTO{
		ID:             rule.ID,
		Title:          rule.Title,
		Description:    rule.Description,
		RecurrenceType: string(rule.RecurrenceType),
		Settings:       rule.Settings,
		ScheduledTime:  rule.ScheduledTime,
		Timezone:       rule.Timezone,
	}
}
