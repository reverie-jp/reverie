package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type DeletePost struct {
	postGateway postgw.Gateway
}

func NewDeletePost(postGateway postgw.Gateway) *DeletePost {
	return &DeletePost{postGateway: postGateway}
}

func (uc *DeletePost) Execute(ctx context.Context, input DeletePostInput, authorID ulid.ULID) error {
	if err := input.Validate(); err != nil {
		return err
	}

	postID, err := ulid.Parse(input.PostID)
	if err != nil {
		return xerrors.ErrInvalidArgument.WithMessage("invalid post_id")
	}

	return uc.postGateway.DeletePost(ctx, postID, authorID)
}
