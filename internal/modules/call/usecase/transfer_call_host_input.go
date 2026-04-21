package usecase

import (
	"errors"
	"strings"

	"reverie.jp/reverie/internal/platform/ulid"
)

type TransferCallHostInput struct {
	RequesterID   ulid.ULID
	CallID        ulid.ULID
	NewHostCustomID string
}

func (i TransferCallHostInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("requester_id is required")
	}
	if i.CallID.IsZero() {
		return errors.New("call_id is required")
	}
	if strings.TrimSpace(i.NewHostCustomID) == "" {
		return errors.New("new_host is required")
	}
	return nil
}
