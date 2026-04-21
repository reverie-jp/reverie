package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type UnbanCallParticipantInput struct {
	RequesterID ulid.ULID
	CallID      ulid.ULID
	UserID      ulid.ULID
}

func (i UnbanCallParticipantInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	if i.CallID.IsZero() {
		return errors.New("call_id is required")
	}
	if i.UserID.IsZero() {
		return errors.New("user_id is required")
	}
	return nil
}
