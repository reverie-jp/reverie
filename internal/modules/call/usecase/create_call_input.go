package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

type CreateCallInput struct {
	RequesterID ulid.ULID
	Visibility  entity.CallVisibility
}

func (i CreateCallInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("authentication required to create a call")
	}
	switch i.Visibility {
	case entity.CallVisibilityOpen, entity.CallVisibilityUsersOnly, entity.CallVisibilityLocked:
		return nil
	default:
		return errors.New("invalid visibility")
	}
}
