package usecase

import (
	"encoding/base64"
	"time"
)

const defaultPageSize int32 = 20
const maxPageSize int32 = 100

// decodePageToken は base64 エンコードされたページトークンをデコードして time.Time を返す。
// 空文字列の場合は nil を返す（最初のページ）。
func decodePageToken(token string) (*time.Time, error) {
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

func normalizePageSize(size int32) int32 {
	if size <= 0 {
		return defaultPageSize
	}
	if size > maxPageSize {
		return maxPageSize
	}
	return size
}
