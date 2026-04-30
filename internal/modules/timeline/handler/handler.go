package handler

import (
	"reverie.jp/reverie/internal/gen/pb/timeline/v1/timelinev1connect"
	timelineusecase "reverie.jp/reverie/internal/modules/timeline/usecase"
)

type Handler struct {
	timelinev1connect.UnimplementedTimelineServiceHandler
	listPublicTimeline    *timelineusecase.ListPublicTimeline
	listFollowingTimeline *timelineusecase.ListFollowingTimeline
}

func New(
	listPublicTimeline *timelineusecase.ListPublicTimeline,
	listFollowingTimeline *timelineusecase.ListFollowingTimeline,
) *Handler {
	return &Handler{
		listPublicTimeline:    listPublicTimeline,
		listFollowingTimeline: listFollowingTimeline,
	}
}
