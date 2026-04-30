package timeline

import (
	"reverie.jp/reverie/internal/gen/pb/timeline/v1/timelinev1connect"
	"reverie.jp/reverie/internal/gen/sqlc"
	postgw "reverie.jp/reverie/internal/modules/post/gateway"
	"reverie.jp/reverie/internal/modules/timeline/handler"
	timelineusecase "reverie.jp/reverie/internal/modules/timeline/usecase"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
)

func InitModule(q sqlc.Querier, userGateway usergw.Gateway) timelinev1connect.TimelineServiceHandler {
	postGateway := postgw.New(q, userGateway)

	listPublicTimeline := timelineusecase.NewListPublicTimeline(postGateway)
	listFollowingTimeline := timelineusecase.NewListFollowingTimeline(postGateway)

	return handler.New(listPublicTimeline, listFollowingTimeline)
}
