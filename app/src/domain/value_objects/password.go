package value_objects

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

type Password string

func NewPassword(raw string) (Password, error) {
	if len(raw) < 8 {
		return "", errors.New("password must be at least 8 characters")
	}
	return Password(raw), nil
}

func (p Password) Hash() string {
	h := sha256.Sum256([]byte(p))
	return hex.EncodeToString(h[:])
}
