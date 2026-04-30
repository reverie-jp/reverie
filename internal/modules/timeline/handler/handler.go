package handler

import (
	"reverie.jp/reverie/internal/gen/pb/timeline/v1/timelinev1connect"
	"reverie.jp/reverie/internal/modules/post/usecase"
)

type Handler struct {
	timelinev1connect.UnimplementedTimelineServiceHandler
	listPublicTimeline    *usecase.ListTimeline
	listFollowingTimeline *usecase.ListFollowingTimeline
}

func New(
	listPublicTimeline *usecase.ListTimeline,
	listFollowingTimeline *usecase.ListFollowingTimeline,
) *Handler {
	return &Handler{
		listPublicTimeline:    listPublicTimeline,
		listFollowingTimeline: listFollowingTimeline,
	}
}
