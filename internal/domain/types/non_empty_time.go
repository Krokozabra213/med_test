package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type NonEmptyTime struct {
	value time.Time
}

func (n NonEmptyTime) Value() time.Time {
	return n.value
}

func NewNonEmptyTime(s string) (NonEmptyTime, error) {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return NonEmptyTime{}, errors.New("time string should not be empty")
	}

	t, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return NonEmptyTime{}, fmt.Errorf("invalid time format (expected RFC3339): %w", err)
	}

	if t.IsZero() {
		return NonEmptyTime{}, errors.New("time is zero")
	}

	return NonEmptyTime{value: t}, nil
}

type NonEmptyString struct {
	value string
}

func (n NonEmptyString) Value() string {
	return n.value
}

func NewNonEmptyString(s string) (NonEmptyString, error) {
	trimmed := strings.TrimSpace(s)

	if len(trimmed) == 0 {
		return NonEmptyString{}, errors.New("string shouldn't be empty")
	}

	return NonEmptyString{
		value: trimmed,
	}, nil
}
