package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListPostRepostsInput struct {
	PostID    string
	PageToken string
	PageSize  int32
}

type ListPostReposts struct {
	postGateway postgw.Gateway
}

func NewListPostReposts(postGateway postgw.Gateway) *ListPostReposts {
	return &ListPostReposts{postGateway: postGateway}
}

func (uc *ListPostReposts) Execute(ctx context.Context, input ListPostRepostsInput, requestorID ulid.ULID) ([]*PostOutput, error) {
	postID, err := ulid.Parse(input.PostID)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid post_id")
	}

	cursor, err := decodePageToken(input.PageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListPostReposts(ctx, postgw.ListPostRepostsParams{
		PostID: postID,
		Cursor: cursor,
		Limit:  normalizePageSize(input.PageSize),
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
