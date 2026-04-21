package usecase

import (
	"errors"
	"strings"

	"reverie.jp/reverie/internal/platform/ulid"
)

type MuteCallParticipantInput struct {
	RequesterID ulid.ULID
	CallID      ulid.ULID
	Identity    string
}

func (i MuteCallParticipantInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	if i.CallID.IsZero() {
		return errors.New("call_id is required")
	}
	if strings.TrimSpace(i.Identity) == "" {
		return errors.New("identity is required")
	}
	return nil
}
