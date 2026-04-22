package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type ListFollowingCallsInput struct {
	RequesterID ulid.ULID
	PageSize    int32
	PageToken   string
}

func (i ListFollowingCallsInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	return nil
}
