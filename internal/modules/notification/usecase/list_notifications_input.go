package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

const (
	defaultListPageSize int32 = 30
	maxListPageSize     int32 = 100
)

type ListNotificationsInput struct {
	RequesterID ulid.ULID
	PageSize    int32
	PageToken   string
}

func (i ListNotificationsInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	return nil
}
