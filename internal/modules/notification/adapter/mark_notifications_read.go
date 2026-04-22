package adapter

import (
	"context"

	"connectrpc.com/connect"

	"reverie.jp/reverie/internal/application/server/interceptor"
	notificationv1 "reverie.jp/reverie/internal/gen/pb/notification/v1"
	"reverie.jp/reverie/internal/modules/notification/usecase"
	"reverie.jp/reverie/internal/platform/resourcename"
	"reverie.jp/reverie/internal/platform/ulid"
	"reverie.jp/reverie/internal/platform/xerrors"
)

func FromMarkNotificationsReadRequest(ctx context.Context, req *connect.Request[notificationv1.MarkNotificationsReadRequest]) (usecase.MarkNotificationsReadInput, error) {
	requesterID, ok := interceptor.UserIDFromContext(ctx)
	if !ok {
		return usecase.MarkNotificationsReadInput{}, xerrors.ErrUnauthenticated
	}
	ids := make([]ulid.ULID, 0, len(req.Msg.Names))
	for _, name := range req.Msg.Names {
		_, id, err := resourcename.ParseNotification(name)
		if err != nil {
			return usecase.MarkNotificationsReadInput{}, err
		}
		ids = append(ids, id)
	}
	return usecase.MarkNotificationsReadInput{
		RequesterID: requesterID,
		IDs:         ids,
	}, nil
}

func ToMarkNotificationsReadResponse(output *usecase.MarkNotificationsReadOutput) *connect.Response[notificationv1.MarkNotificationsReadResponse] {
	return connect.NewResponse(&notificationv1.MarkNotificationsReadResponse{
		MarkedCount: output.MarkedCount,
	})
}
