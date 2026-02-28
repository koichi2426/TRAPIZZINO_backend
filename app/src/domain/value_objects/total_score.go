package value_objects

import "errors"

// TotalScore は、最終的な「運命の1軒」を決定するための統合計算スコアです。
type TotalScore float64

func NewTotalScore(value float64) (TotalScore, error) {
	if value < 0 {
		return 0, errors.New("total score cannot be negative")
	}
	return TotalScore(value), nil
}

func (s TotalScore) Float64() float64 {
	return float64(s)
}