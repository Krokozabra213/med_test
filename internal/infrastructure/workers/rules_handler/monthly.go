package ruleshandler

import (
	"bytes"
	"encoding/json"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type MonthlyDayHandler struct{}

func (h *MonthlyDayHandler) CanHandle(rt taskdomain.RecurrenceType) bool {
	return rt == taskdomain.RecurrenceMonthlyDay
}

func (h *MonthlyDayHandler) Handle(rule *taskdomain.Rule, date time.Time) (bool, error) {

	dec := json.NewDecoder(bytes.NewReader(rule.Settings))
	dec.DisallowUnknownFields()

	var s taskdomain.MonthlyDaySettings

	if err := dec.Decode(&s); err != nil {
		return false, err
	}

	return date.Day() == s.Day, nil
}
