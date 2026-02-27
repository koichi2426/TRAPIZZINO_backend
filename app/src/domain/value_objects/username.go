package value_objects

import (
	"errors"
	"unicode/utf8"
)

type Username string

func NewUsername(value string) (Username, error) {
	if utf8.RuneCountInString(value) < 3 || utf8.RuneCountInString(value) > 32 {
		return "", errors.New("username must be 3-32 chars")
	}
	return Username(value), nil
}

func (u Username) String() string {
	return string(u)
}
