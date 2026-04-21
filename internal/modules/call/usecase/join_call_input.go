package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/validation"
)

type JoinCallInput struct {
	RoomID string `validate:"required"`
	// Zero value means the caller is a guest.
	RequesterID      ulid.ULID
	GuestDisplayName string
}

func (i JoinCallInput) Validate() error {
	if err := validation.CheckStruct(i); err != nil {
		return err
	}
	if i.RequesterID.IsZero() && i.GuestDisplayName == "" {
		return errors.New("guest_display_name is required for guest callers")
	}
	return nil
}
