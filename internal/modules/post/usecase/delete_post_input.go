package usecase

import "reverie.jp/reverie/internal/platform/validation"

type DeletePostInput struct {
	PostID string `validate:"required"`
}

func (i DeletePostInput) Validate() error {
	return validation.CheckStruct(i)
}
