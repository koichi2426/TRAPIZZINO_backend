package value_objects

import (
	"errors"
	"regexp"
)

type Email string

var emailRegex = regexp.MustCompile(`^[\w._%+-]+@[\w.-]+\\.[a-zA-Z]{2,}$`)

func NewEmail(value string) (Email, error) {
	if !emailRegex.MatchString(value) {
		return "", errors.New("invalid email format")
	}
	return Email(value), nil
}

func (e Email) String() string {
	return string(e)
}
