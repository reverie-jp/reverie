package usecase

import (
	"context"

	"reverie.jp/reverie/internal/platform/pagetoken"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	postusecase "reverie.jp/reverie/internal/modules/post/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListFollowingTimeline struct {
	postGateway postgw.Gateway
}

func NewListFollowingTimeline(postGateway postgw.Gateway) *ListFollowingTimeline {
	return &ListFollowingTimeline{postGateway: postGateway}
}

func (uc *ListFollowingTimeline) Execute(ctx context.Context, input ListTimelineInput, requestorID ulid.ULID) ([]*postusecase.PostOutput, error) {
	cursor, err := pagetoken.Decode(input.PageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListFollowingTimeline(ctx, postgw.ListFollowingTimelineParams{
		FollowerID: requestorID,
		Cursor:     cursor,
		Limit:      pagetoken.NormalizePageSize(input.PageSize),
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
