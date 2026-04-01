package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/validation"
)

type DeleteAccountInput struct {
	UserID          ulid.ULID
	ConfirmCustomID string `validate:"required"`
}

func (i DeleteAccountInput) Validate() error {
	if i.UserID.IsZero() {
		return errors.New("user_id is required")
	}
	return validation.CheckStruct(i)
}
