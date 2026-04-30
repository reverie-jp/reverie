package usecase

import (
	"context"

	"reverie.jp/reverie/internal/platform/pagetoken"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	postusecase "reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListPublicTimeline struct {
	postGateway postgw.Gateway
}

func NewListPublicTimeline(postGateway postgw.Gateway) *ListPublicTimeline {
	return &ListPublicTimeline{postGateway: postGateway}
}

func (uc *ListPublicTimeline) Execute(ctx context.Context, input ListTimelineInput, requestorID ulid.ULID) ([]*postusecase.PostOutput, error) {
	cursor, err := pagetoken.Decode(input.PageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListTimeline(ctx, postgw.ListTimelineParams{
		Cursor: cursor,
		Limit:  pagetoken.NormalizePageSize(input.PageSize),
	}, requestorID)
	if err != nil {
		return nil, err
	}

	outputs := make([]*postusecase.PostOutput, len(views))
	for i, v := range views {
		outputs[i] = postusecase.ToPostOutput(v)
	}
	return outputs, nil
}
