package value_objects

import "errors"

type ResonanceCount int

func NewResonanceCount(value int) (ResonanceCount, error) {
	if value < 0 {
		return 0, errors.New("resonance count cannot be negative")
	}
	return ResonanceCount(value), nil
}

func (r ResonanceCount) Int() int {
	return int(r)
}