package usecase

import (
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/validation"
)

type GetCallInput struct {
	RequesterID ulid.ULID
	CallID      ulid.ULID `validate:"required"`
	// ViewerIdentity is the LiveKit identity of the caller, if they are
	// actively connected. For authenticated callers the adapter fills this
	// with "user:<ulid>". For guests the adapter relays the client-provided
	// guest_identity. Used to grant LOCKED visibility to currently-connected
	// participants.
	ViewerIdentity string
}

func (i GetCallInput) Validate() error {
	return validation.CheckStruct(i)
}
