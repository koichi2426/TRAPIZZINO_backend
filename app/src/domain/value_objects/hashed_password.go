package value_objects

import "errors"

type HashedPassword string

func NewHashedPassword(value string) (HashedPassword, error) {
	if len(value) < 8 {
		return "", errors.New("hashed password too short")
	}
	return HashedPassword(value), nil
}

func (h HashedPassword) String() string {
	return string(h)
}
