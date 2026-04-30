package usecase

import "reverie.jp/reverie/internal/platform/xerrors"

type ListUserLikedPostsInput struct {
	UserID    string
	PageToken string
	PageSize  int32
}

func (in ListUserLikedPostsInput) Validate() error {
	if in.UserID == "" {
		return xerrors.ErrInvalidArgument.WithMessage("user_id is required")
	}
	return nil
}
