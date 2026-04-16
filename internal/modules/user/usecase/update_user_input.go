package usecase

import (
	"time"

	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/validation"
)

type UpdateUserInput struct {
	DisplayName string  `validate:"required,min=1,max=20"`
	Biography   string `validate:"max=160"`
	IsPrivate   bool
	Birthdate   *time.Time
}

func (i UpdateUserInput) Validate() error {
	return validation.CheckStruct(i)
}

type UpdateUserInputParsed struct {
	UserID ulid.ULID
}
