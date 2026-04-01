package usecase

import "reverie.jp/reverie/internal/platform/validation"

type SocialLoginInput struct {
	Provider string `validate:"required,oneof=google"`
	Code     string `validate:"required"`
}

func (i SocialLoginInput) Validate() error {
	return validation.CheckStruct(i)
}
