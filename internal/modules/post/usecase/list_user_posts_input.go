package usecase

import "reverie.jp/reverie/internal/platform/validation"

type ListUserPostsInput struct {
	UserID    string `validate:"required"`
	PageToken string
	PageSize  int32
}

func (i ListUserPostsInput) Validate() error {
	return validation.CheckStruct(i)
}
