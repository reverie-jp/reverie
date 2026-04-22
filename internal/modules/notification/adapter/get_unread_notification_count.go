package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	notificationv1 "reverie.jp/reverie/internal/gen/pb/notification/v1"
	"reverie.jp/reverie/internal/modules/notification/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromGetUnreadNotificationCountRequest(ctx context.Context, _ *connect.Request[notificationv1.GetUnreadNotificationCountRequest]) (usecase.GetUnreadNotificationCountInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.GetUnreadNotificationCountInput{}, xerrors.ErrUnauthenticated
	}
	return usecase.GetUnreadNotificationCountInput{RequesterID: requesterID}, nil
}

func ToGetUnreadNotificationCountResponse(output *usecase.GetUnreadNotificationCountOutput) *connect.Response[notificationv1.GetUnreadNotificationCountResponse] {
	return connect.NewResponse(&notificationv1.GetUnreadNotificationCountResponse{
		Count: output.Count,
	})
}
