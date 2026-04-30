package usecase

import (
	"context"

	"reverie.jp/reverie/internal/platform/pagetoken"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListUserLikedPosts struct {
	postGateway postgw.Gateway
}

func NewListUserLikedPosts(postGateway postgw.Gateway) *ListUserLikedPosts {
	return &ListUserLikedPosts{postGateway: postGateway}
}

func (uc *ListUserLikedPosts) Execute(ctx context.Context, input ListUserLikedPostsInput, requestorID ulid.ULID) ([]*PostOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	targetID, err := ulid.Parse(input.UserID)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid user_id")
	}

	cursor, err := pagetoken.Decode(input.PageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListUserLikedPosts(ctx, postgw.ListUserLikedPostsParams{
		UserID: targetID,
		Cursor: cursor,
		Limit:  pagetoken.NormalizePageSize(input.PageSize),
	}, requestorID)
	if err != nil {
		return nil, err
	}

	outputs := make([]*PostOutput, len(views))
	for i, v := range views {
		outputs[i] = ToPostOutput(v)
	}
	return outputs, nil
}
