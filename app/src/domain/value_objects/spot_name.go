package value_objects

import (
	"errors"
	"unicode/utf8"
)

type SpotName string

func NewSpotName(value string) (SpotName, error) {
	if utf8.RuneCountInString(value) < 1 || utf8.RuneCountInString(value) > 64 {
		return "", errors.New("spot name must be 1-64 chars")
	}
	return SpotName(value), nil
}

func (n SpotName) String() string {
	return string(n)
}
