package handler

import (
	"context"

	"connectrpc.com/connect"

	timelinev1 "reverie.jp/reverie/internal/gen/pb/timeline/v1"
	timelineadapter "reverie.jp/reverie/internal/modules/timeline/adapter"
)

func (h *Handler) ListPublicTimeline(ctx context.Context, req *connect.Request[timelinev1.ListPublicTimelineRequest]) (*connect.Response[timelinev1.ListPublicTimelineResponse], error) {
	input, userID, err := timelineadapter.FromListPublicTimelineRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	outputs, err := h.listPublicTimeline.Execute(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	return timelineadapter.ToListPublicTimelineResponse(outputs), nil
}
