package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type LikePostInput struct {
	PostID string `validate:"required"`
}

type UnlikePostInput struct {
	PostID string `validate:"required"`
}

type LikePost struct {
	postGateway postgw.Gateway
}

type UnlikePost struct {
	postGateway postgw.Gateway
}

func NewLikePost(postGateway postgw.Gateway) *LikePost {
	return &LikePost{postGateway: postGateway}
}

func NewUnlikePost(postGateway postgw.Gateway) *UnlikePost {
	return &UnlikePost{postGateway: postGateway}
}

func (uc *LikePost) Execute(ctx context.Context, input LikePostInput, userID ulid.ULID) (*PostOutput, error) {
	postID, err := ulid.Parse(input.PostID)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid post_id")
	}

	view, err := uc.postGateway.LikePost(ctx, postID, userID)
	if err != nil {
		return nil, err
	}

	return toPostOutput(view), nil
}

func (uc *UnlikePost) Execute(ctx context.Context, input UnlikePostInput, userID ulid.ULID) (*PostOutput, error) {
	postID, err := ulid.Parse(input.PostID)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid post_id")
	}

	view, err := uc.postGateway.UnlikePost(ctx, postID, userID)
	if err != nil {
		return nil, err
	}

	return toPostOutput(view), nil
}
