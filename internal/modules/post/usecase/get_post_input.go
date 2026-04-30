package usecase

import "reverie.jp/reverie/internal/platform/validation"

type GetPostInput struct {
	AuthorCustomID string `validate:"required"`
	ShortID        string `validate:"required"`
}

func (i GetPostInput) Validate() error {
	return validation.CheckStruct(i)
}
