package task

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrRuleNotFound = errors.New("rule not found")
)
