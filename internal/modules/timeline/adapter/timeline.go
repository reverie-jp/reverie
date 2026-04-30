package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	timelinev1 "reverie.jp/reverie/internal/gen/pb/timeline/v1"
	postadapter "reverie.jp/reverie/internal/modules/post/adapter"
	postusecase "reverie.jp/reverie/internal/modules/post/usecase"
	timelineusecase "reverie.jp/reverie/internal/modules/timeline/usecase"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromListPublicTimelineRequest(ctx context.Context, req *connect.Request[timelinev1.ListPublicTimelineRequest]) (timelineusecase.ListTimelineInput, ulid.ULID, error) {
	userID, _ := interceptor.UserIDFromContext(ctx)
	return timelineusecase.ListTimelineInput{
		PageToken: req.Msg.PageToken,
		PageSize:  req.Msg.PageSize,
	}, userID, nil
}

func ToListPublicTimelineResponse(outputs []*postusecase.PostOutput) *connect.Response[timelinev1.ListPublicTimelineResponse] {
	return connect.NewResponse(&timelinev1.ListPublicTimelineResponse{
		Posts:         postadapter.ToPosts(outputs),
		NextPageToken: postadapter.NextPageToken(outputs),
	})
}

func FromListFollowingTimelineRequest(ctx context.Context, req *connect.Request[timelinev1.ListFollowingTimelineRequest]) (timelineusecase.ListTimelineInput, ulid.ULID, error) {
	userID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return timelineusecase.ListTimelineInput{}, ulid.ULID{}, xerrors.ErrUnauthenticated
	}
	return timelineusecase.ListTimelineInput{
		PageToken: req.Msg.PageToken,
		PageSize:  req.Msg.PageSize,
	}, userID, nil
}

func ToListFollowingTimelineResponse(outputs []*postusecase.PostOutput) *connect.Response[timelinev1.ListFollowingTimelineResponse] {
	return connect.NewResponse(&timelinev1.ListFollowingTimelineResponse{
		Posts:         postadapter.ToPosts(outputs),
		NextPageToken: postadapter.NextPageToken(outputs),
	})
}
