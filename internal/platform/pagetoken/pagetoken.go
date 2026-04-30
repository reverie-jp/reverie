package pagetoken

import (
	"encoding/base64"
	"time"
)

const DefaultPageSize int32 = 20
const MaxPageSize int32 = 100

func Decode(token string) (*time.Time, error) {
	if token == "" {
		return nil, nil
	}
	raw, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		return nil, err
	}
	t, err := time.Parse("2006-01-02T15:04:05.999999999Z", string(raw))
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func NormalizePageSize(size int32) int32 {
	if size <= 0 {
		return DefaultPageSize
	}
	if size > MaxPageSize {
		return MaxPageSize
	}
	return size
}
