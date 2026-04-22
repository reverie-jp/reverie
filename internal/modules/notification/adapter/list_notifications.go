package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	notificationv1 "reverie.jp/reverie/internal/gen/pb/notification/v1"
	notificationgw "reverie.jp/reverie/internal/modules/notification/gateway"
	"reverie.jp/reverie/internal/modules/notification/usecase"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromListNotificationsRequest(ctx context.Context, req *connect.Request[notificationv1.ListNotificationsRequest]) (usecase.ListNotificationsInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.ListNotificationsInput{}, xerrors.ErrUnauthenticated
	}
	return usecase.ListNotificationsInput{
		RequesterID: requesterID,
		PageSize:    req.Msg.PageSize,
		PageToken:   req.Msg.PageToken,
	}, nil
}

func ToListNotificationsResponse(output *usecase.ListNotificationsOutput) *connect.Response[notificationv1.ListNotificationsResponse] {
	out := make([]*notificationv1.Notification, 0, len(output.Views))
	for _, v := range output.Views {
		if pb := notificationgw.ToProto(v); pb != nil {
			out = append(out, pb)
		}
	}
	return connect.NewResponse(&notificationv1.ListNotificationsResponse{
		Notifications: out,
		NextPageToken: output.NextPageToken,
	})
}
