package value_objects

import "errors"

type Reason string

func NewReason(value string) (Reason, error) {
	if value == "" {
		return "", errors.New("reason cannot be empty")
	}
	return Reason(value), nil
}

func (r Reason) String() string {
	return string(r)
}