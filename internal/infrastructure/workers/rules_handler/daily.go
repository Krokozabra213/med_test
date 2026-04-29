package ruleshandler

import (
	"bytes"
	"encoding/json"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type DailyHandler struct{}

func (h *DailyHandler) CanHandle(rt taskdomain.RecurrenceType) bool {
	return rt == taskdomain.RecurrenceDaily
}

func (h *DailyHandler) Handle(rule *taskdomain.Rule, date time.Time) (bool, error) {

	dec := json.NewDecoder(bytes.NewReader(rule.Settings))
	dec.DisallowUnknownFields()

	var s taskdomain.DailySettings

	if err := dec.Decode(&s); err != nil {
		return false, err
	}

	daysPassed := int(date.Sub(rule.CreatedAt).Hours() / 24)

	return daysPassed%s.Interval == 0, nil
}
