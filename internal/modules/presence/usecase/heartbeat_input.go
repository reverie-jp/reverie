package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type HeartbeatInput struct {
	RequesterID ulid.ULID
}

func (i HeartbeatInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	return nil
}
