package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type GetUserInput struct {
	RequesterID ulid.ULID
	TargetID    ulid.ULID
}

func (i GetUserInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	if i.TargetID.IsZero() {
		return errors.New("target_id is required")
	}
	return nil
}
