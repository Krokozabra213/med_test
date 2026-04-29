package ruleshandler

import (
	"bytes"
	"encoding/json"
	"time"

	taskdomain "example.com/taskservice/internal/domain/task"
)

type SpecificDatesHandler struct{}

func (h *SpecificDatesHandler) CanHandle(rt taskdomain.RecurrenceType) bool {
	return rt == taskdomain.RecurrenceSpecific
}

func (h *SpecificDatesHandler) Handle(rule *taskdomain.Rule, date time.Time) (bool, error) {

	dec := json.NewDecoder(bytes.NewReader(rule.Settings))
	dec.DisallowUnknownFields()

	var s taskdomain.SpecificDatesSettings

	if err := dec.Decode(&s); err != nil {
		return false, err
	}

	today := date.Format("2006-01-02")

	for _, d := range s.Dates {
		if d == today {
			return true, nil
		}
	}

	return false, nil
}
