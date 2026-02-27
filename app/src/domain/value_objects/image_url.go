package value_objects

import (
	"errors"
	"net/url"
)

type ImageURL string

func NewImageURL(value string) (ImageURL, error) {
	u, err := url.ParseRequestURI(value)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "", errors.New("invalid image url")
	}
	return ImageURL(value), nil
}

func (u ImageURL) String() string {
	return string(u)
}
