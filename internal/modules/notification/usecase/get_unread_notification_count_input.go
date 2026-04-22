package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type GetUnreadNotificationCountInput struct {
	RequesterID ulid.ULID
}

func (i GetUnreadNotificationCountInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	return nil
}
