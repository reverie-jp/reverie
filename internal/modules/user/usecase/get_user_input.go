package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

// GetUserInput identifies the target user by either ULID (self-resolution path,
// used by GetMyUser) or custom_id (the public path from URL, used by GetUser).
// Exactly one of TargetID / TargetCustomID should be set.
type GetUserInput struct {
	RequesterID    ulid.ULID
	TargetID       ulid.ULID
	TargetCustomID string
}

func (i GetUserInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	if i.TargetID.IsZero() && i.TargetCustomID == "" {
		return errors.New("target_id or target_custom_id is required")
	}
	return nil
}
