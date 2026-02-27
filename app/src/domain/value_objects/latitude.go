package value_objects

import "errors"

type Latitude float64

func NewLatitude(value float64) (Latitude, error) {
	if value < -90 || value > 90 {
		return 0, errors.New("latitude out of range")
	}
	return Latitude(value), nil
}

func (lat Latitude) Value() float64 {
	return float64(lat)
}
