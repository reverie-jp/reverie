package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListPostRepliesInput struct {
	PostID    string
	PageToken string
	PageSize  int32
}

type ListPostReplies struct {
	postGateway postgw.Gateway
}

func NewListPostReplies(postGateway postgw.Gateway) *ListPostReplies {
	return &ListPostReplies{postGateway: postGateway}
}

func (uc *ListPostReplies) Execute(ctx context.Context, input ListPostRepliesInput, requestorID ulid.ULID) ([]*PostOutput, error) {
	postID, err := ulid.Parse(input.PostID)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid post_id")
	}

	cursor, err := decodePageToken(input.PageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListPostReplies(ctx, postgw.ListPostRepliesParams{
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
