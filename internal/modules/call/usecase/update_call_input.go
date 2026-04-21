package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

type UpdateCallInput struct {
	RequesterID ulid.ULID
	CallID      ulid.ULID
	Visibility  entity.CallVisibility
	UpdateMask  []string
}

func (i UpdateCallInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("authentication required")
	}
	if i.CallID.IsZero() {
		return errors.New("call_id is required")
	}
	if len(i.UpdateMask) == 0 {
		return errors.New("update_mask is required")
	}
	for _, path := range i.UpdateMask {
		if path != "visibility" {
			return errors.New("only visibility is updatable")
		}
	}
	switch i.Visibility {
	case entity.CallVisibilityOpen, entity.CallVisibilityUsersOnly, entity.CallVisibilityLocked:
		return nil
	default:
		return errors.New("invalid visibility")
	}
}
