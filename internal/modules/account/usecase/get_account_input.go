package usecase

import (
	"errors"

	"reverie.jp/reverie/internal/platform/ulid"
)

type GetAccountInput struct {
	UserID ulid.ULID
}

func (i GetAccountInput) Validate() error {
	if i.UserID.IsZero() {
		return errors.New("user_id is required")
	}
	return nil
}
