package usecase

import "reverie.jp/reverie/internal/platform/validation"

type RefreshTokenInput struct {
	RefreshToken string `validate:"required"`
}

func (i RefreshTokenInput) Validate() error {
	return validation.CheckStruct(i)
}
