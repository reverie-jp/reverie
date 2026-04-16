package usecase

import (
	"context"

	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

type ListTimeline struct {
	postGateway postgw.Gateway
}

func NewListTimeline(postGateway postgw.Gateway) *ListTimeline {
	return &ListTimeline{postGateway: postGateway}
}

func (uc *ListTimeline) Execute(ctx context.Context, input ListTimelineInput, requestorID ulid.ULID) ([]*PostOutput, error) {
	cursor, err := decodePageToken(input.PageToken)
	if err != nil {
		return nil, xerrors.ErrInvalidArgument.WithMessage("invalid page_token")
	}

	views, err := uc.postGateway.ListTimeline(ctx, postgw.ListTimelineParams{
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
