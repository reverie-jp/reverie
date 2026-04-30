package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListFollowingTimeline struct {
	postGateway postgw.Gateway
}

func NewListFollowingTimeline(postGateway postgw.Gateway) *ListFollowingTimeline {
	return &ListFollowingTimeline{postGateway: postGateway}
}

func (uc *ListFollowingTimeline) Execute(ctx context.Context, input ListTimelineInput, requestorID ulid.ULID) ([]*PostOutput, error) {
	cursor, err := decodePageToken(input.PageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListFollowingTimeline(ctx, postgw.ListFollowingTimelineParams{
		FollowerID: requestorID,
		Cursor:     cursor,
		Limit:      normalizePageSize(input.PageSize),
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
