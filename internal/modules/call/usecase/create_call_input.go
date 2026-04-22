package usecase

import (
	"errors"
	"unicode/utf8"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

const CallTitleMaxLength = 100

type CreateCallInput struct {
	RequesterID ulid.ULID
	Visibility  entity.CallVisibility
	Title       string
}

func (i CreateCallInput) Validate() error {
	if i.RequesterID.IsZero() {
		return errors.New("authentication required to create a call")
	}
	switch i.Visibility {
	case entity.CallVisibilityOpen, entity.CallVisibilityUsersOnly, entity.CallVisibilityLocked:
	default:
		return errors.New("invalid visibility")
	}
	if utf8.RuneCountInString(i.Title) > CallTitleMaxLength {
		return errors.New("title too long")
	}
	return nil
}
