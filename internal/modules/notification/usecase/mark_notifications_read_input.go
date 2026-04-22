package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type MarkNotificationsReadInput struct {
	RequesterID ulid.ULID
	IDs         []ulid.ULID
}

func (i MarkNotificationsReadInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	return nil
}
