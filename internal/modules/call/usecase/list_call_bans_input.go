package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type ListCallBansInput struct {
	RequesterID ulid.ULID
	CallID      ulid.ULID
	PageSize    int32
	PageToken   string
}

func (i ListCallBansInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	if i.CallID.IsZero() {
		return errors.New("call_id is required")
	}
	return nil
}
