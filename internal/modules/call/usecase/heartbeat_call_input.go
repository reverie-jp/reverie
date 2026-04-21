package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type HeartbeatCallInput struct {
	RequesterID ulid.ULID
	CallID      ulid.ULID
	Identity    string
}

func (i HeartbeatCallInput) Validate() error {
	if i.CallID.IsZero() {
		return errors.New("call_id is required")
	}
	if i.Identity == "" {
		return errors.New("identity is required")
	}
	return nil
}
