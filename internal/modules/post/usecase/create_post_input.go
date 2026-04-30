package usecase

import (
	"reverie.jp/reverie/internal/platform/validation"
)

type CreatePostInput struct {
	Text      string  `validate:"required,min=1,max=500"`
	ReplyToID *string `validate:"omitempty"`
	RepostID  *string `validate:"omitempty"`
}

func (i CreatePostInput) Validate() error {
	return validation.CheckStruct(i)
}
