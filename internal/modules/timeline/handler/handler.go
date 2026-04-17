package handler

import (
	timelinev1connect "reverie.jp/reverie/internal/gen/pb/timeline/v1/timelinev1connect"
	"reverie.jp/reverie/internal/modules/post/usecase"
)

type Handler struct {
	timelinev1connect.UnimplementedTimelineServiceHandler
	listTimeline          *usecase.ListTimeline
	listFollowingTimeline *usecase.ListFollowingTimeline
}

func New(listTimeline *usecase.ListTimeline, listFollowingTimeline *usecase.ListFollowingTimeline) *Handler {
	return &Handler{
		listTimeline:          listTimeline,
		listFollowingTimeline: listFollowingTimeline,
	}
}
