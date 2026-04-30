package usecase

import (
	"context"

	"reverie.jp/reverie/internal/platform/pagetoken"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListPostRepostsInput struct {
	AuthorCustomID string
	ShortID        string
	PageToken      string
	PageSize       int32
}

type ListPostReposts struct {
	postGateway postgw.Gateway
}

func NewListPostReposts(postGateway postgw.Gateway) *ListPostReposts {
	return &ListPostReposts{postGateway: postGateway}
}

func (uc *ListPostReposts) Execute(ctx context.Context, input ListPostRepostsInput, requestorID ulid.ULID) ([]*PostOutput, error) {
	cursor, err := pagetoken.Decode(input.PageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListPostReposts(ctx, input.AuthorCustomID, input.ShortID, postgw.ListPostRepostsParams{
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
