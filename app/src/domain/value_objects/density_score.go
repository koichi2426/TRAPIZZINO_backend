package value_objects

import "errors"

// DensityScore は、特定のメッシュに対する注目度（登録・上書きの延べ回数）を表します。
type DensityScore int

func NewDensityScore(value int) (DensityScore, error) {
	if value < 0 {
		return 0, errors.New("density score cannot be negative")
	}
	return DensityScore(value), nil
}

func (d DensityScore) Int() int {
	return int(d)
}