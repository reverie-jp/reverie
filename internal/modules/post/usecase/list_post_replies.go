package usecase

import (
	"context"

	"reverie.jp/reverie/internal/platform/pagetoken"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListPostRepliesInput struct {
	AuthorCustomID string
	ShortID        string
	PageToken      string
	PageSize       int32
}

type ListPostReplies struct {
	postGateway postgw.Gateway
}

func NewListPostReplies(postGateway postgw.Gateway) *ListPostReplies {
	return &ListPostReplies{postGateway: postGateway}
}

func (uc *ListPostReplies) Execute(ctx context.Context, input ListPostRepliesInput, requestorID ulid.ULID) ([]*PostOutput, error) {
	cursor, err := pagetoken.Decode(input.PageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListPostReplies(ctx, input.AuthorCustomID, input.ShortID, postgw.ListPostRepliesParams{
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
