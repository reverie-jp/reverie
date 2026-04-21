package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type JoinCallInput struct {
	CallID ulid.ULID
	// Zero value means the caller is a guest.
	RequesterID      ulid.ULID
	GuestDisplayName string
}

func (i JoinCallInput) Validate() error {
	if i.CallID.IsZero() {
		return errors.New("call_id is required")
	}
	if i.RequesterID.IsZero() && i.GuestDisplayName == "" {
		return errors.New("guest_display_name is required for guest callers")
	}
	return nil
}
