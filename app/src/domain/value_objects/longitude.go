package value_objects

import "errors"

type Longitude float64

func NewLongitude(value float64) (Longitude, error) {
	if value < -180 || value > 180 {
		return 0, errors.New("longitude out of range")
	}
	return Longitude(value), nil
}

func (lng Longitude) Value() float64 {
	return float64(lng)
}
