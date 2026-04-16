package usecase

import (
	"context"
	"time"

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
	params := postgw.ListTimelineParams{
		Limit: input.Limit,
	}

	if input.Cursor != nil {
		t, err := time.Parse(time.RFC3339Nano, *input.Cursor)
		if err != nil {
			return nil, xerrors.ErrInvalidArgument.WithMessage("invalid cursor")
		}
		params.Cursor = &t
	}

	views, err := uc.postGateway.ListTimeline(ctx, params, requestorID)
	if err != nil {
		return nil, err
	}

	outputs := make([]*PostOutput, len(views))
	for i, v := range views {
		outputs[i] = toPostOutput(v)
	}
	return outputs, nil
}
