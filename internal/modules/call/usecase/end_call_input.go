package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type EndCallInput struct {
	RequesterID ulid.ULID
	CallID      ulid.ULID
}

func (i EndCallInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	if i.CallID.IsZero() {
		return errors.New("call_id is required")
	}
	return nil
}
