package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type GetPost struct {
	postGateway postgw.Gateway
}

func NewGetPost(postGateway postgw.Gateway) *GetPost {
	return &GetPost{postGateway: postGateway}
}

func (uc *GetPost) Execute(ctx context.Context, input GetPostInput, requestorID ulid.ULID) (*PostOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	postID, err := ulid.Parse(input.PostID)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid post_id")
	}

	view, err := uc.postGateway.GetPost(ctx, postID, requestorID)
	if err != nil {
		return nil, err
	}
	if view == nil {
		return nil, xerrors.ErrPostNotFound
	}

	return toPostOutput(view), nil
}
