package value_objects

import (
    "errors"
    "net/url"
)

type ImageURL string

func NewImageURL(value string) (ImageURL, error) {
    // 空文字（JSONでnullの場合も含む）は有効な「画像なし」として扱う
    if value == "" {
        return "", nil
    }

    u, err := url.ParseRequestURI(value)
    if err != nil || u.Scheme == "" || u.Host == "" {
        return "", errors.New("invalid image url")
    }
    return ImageURL(value), nil
}

func (u ImageURL) String() string {
    return string(u)
}