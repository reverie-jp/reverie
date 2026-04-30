package handler

import (
	"context"

	"connectrpc.com/connect"

	timelinev1 "reverie.jp/reverie/internal/gen/pb/timeline/v1"
	timelineadapter "reverie.jp/reverie/internal/modules/timeline/adapter"
)

func (h *Handler) ListFollowingTimeline(ctx context.Context, req *connect.Request[timelinev1.ListFollowingTimelineRequest]) (*connect.Response[timelinev1.ListFollowingTimelineResponse], error) {
	input, userID, err := timelineadapter.FromListFollowingTimelineRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	outputs, err := h.listFollowingTimeline.Execute(ctx, input, userID)
	if err != nil {
		return nil, err
	}
	return timelineadapter.ToListFollowingTimelineResponse(outputs), nil
}
