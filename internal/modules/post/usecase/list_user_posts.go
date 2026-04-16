package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListUserPosts struct {
	postGateway postgw.Gateway
}

func NewListUserPosts(postGateway postgw.Gateway) *ListUserPosts {
	return &ListUserPosts{postGateway: postGateway}
}

func (uc *ListUserPosts) Execute(ctx context.Context, input ListUserPostsInput, requestorID ulid.ULID) ([]*PostOutput, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	targetID, err := ulid.Parse(input.UserID)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid user_id")
	}

	cursor, err := decodePageToken(input.PageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListUserPosts(ctx, postgw.ListUserPostsParams{
		AuthorID: targetID,
		Cursor:   cursor,
		Limit:    normalizePageSize(input.PageSize),
	}, requestorID)
	if err != nil {
		return nil, err
	}

	outputs := make([]*PostOutput, len(views))
	for i, v := range views {
		outputs[i] = toPostOutput(v)
	}
	return outputs, nil
}
