package usecase

import (
	"context"

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

func (uc *ListUserLikedPosts) Execute(ctx context.Context, userIDStr string, pageToken string, pageSize int32, requestorID ulid.ULID) ([]*PostOutput, error) {
	targetID, err := ulid.Parse(userIDStr)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid user_id")
	}

	cursor, err := decodePageToken(pageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListUserLikedPosts(ctx, postgw.ListUserLikedPostsParams{
		UserID: targetID,
		Cursor: cursor,
		Limit:  normalizePageSize(pageSize),
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
