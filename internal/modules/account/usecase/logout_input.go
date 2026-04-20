package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/validation"
)

type LogoutInput struct {
	UserID       ulid.ULID
	RefreshToken string `validate:"required"`
}

func (i LogoutInput) Validate() error {
	if i.UserID.IsZero() {
		return errors.New("user_id is required")
	}
	return validation.CheckStruct(i)
}
