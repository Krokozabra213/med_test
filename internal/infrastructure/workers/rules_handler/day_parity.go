package ruleshandler

import (
	"bytes"
	"encoding/json"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type DayParityHandler struct{}

func (h *DayParityHandler) CanHandle(rt taskdomain.RecurrenceType) bool {
	return rt == taskdomain.RecurrenceDayParity
}

func (h *DayParityHandler) Handle(rule *taskdomain.Rule, date time.Time) (bool, error) {

	dec := json.NewDecoder(bytes.NewReader(rule.Settings))
	dec.DisallowUnknownFields()

	var s taskdomain.DayParitySettings

	if err := dec.Decode(&s); err != nil {
		return false, err
	}

	if s.Parity == "even" {
		return date.Day()%2 == 0, nil
	}
	return date.Day()%2 != 0, nil
}
