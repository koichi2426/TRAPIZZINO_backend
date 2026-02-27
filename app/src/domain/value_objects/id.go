package value_objects

import "errors"

type ID int

func NewID(value int) (ID, error) {
	if value < 0 {
		return 0, errors.New("ID must be non-negative")
	}
	return ID(value), nil
}

func (id ID) Value() int {
	return int(id)
}
