package value_objects

import (
	"errors"
	"unicode/utf8"
)

type Caption string

func NewCaption(value string) (Caption, error) {
	if utf8.RuneCountInString(value) > 256 {
		return "", errors.New("caption must be <= 256 chars")
	}
	return Caption(value), nil
}

func (c Caption) String() string {
	return string(c)
}
