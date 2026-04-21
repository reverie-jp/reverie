package usecase

import (
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/validation"
)

type GetUserParticipatingCallInput struct {
	RequesterID  ulid.ULID
	UserCustomID string `validate:"required"`
}

func (i GetUserParticipatingCallInput) Validate() error {
	return validation.CheckStruct(i)
}
