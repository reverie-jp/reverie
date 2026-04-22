package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type UnfollowUserInput struct {
	RequesterID    ulid.ULID
	TargetCustomID string
}

func (i UnfollowUserInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	if i.TargetCustomID == "" {
		return errors.New("target_custom_id is required")
	}
	return nil
}
