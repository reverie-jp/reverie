package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
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
	return uc.postGateway.DeletePost(ctx, input.AuthorCustomID, input.ShortID, authorID)
}
