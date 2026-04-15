package usecase

import (
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/validation"
)

type GetUserInput struct {
	UserID string `validate:"required"`
}

func (i GetUserInput) Validate() error {
	return validation.CheckStruct(i)
}

type GetUserInputParsed struct {
	UserID ulid.ULID
}
