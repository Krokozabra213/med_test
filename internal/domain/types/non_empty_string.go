package domainTypes

import (
	"errors"
	"strings"
)

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
