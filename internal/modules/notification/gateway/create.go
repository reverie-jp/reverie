package gateway

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/timestamppb"

	eventv1 "reverie.jp/reverie/internal/gen/pb/event/v1"
	"reverie.jp/reverie/internal/modules/notification/repository"
	"reverie.jp/reverie/internal/platform/events"
	"reverie.jp/reverie/internal/platform/ulid"
)

func (g *gatewayImpl) Create(ctx context.Context, params CreateParams) (*NotificationView, error) {
	if params.ActorUserID != nil && *params.ActorUserID == params.RecipientUserID {
		return nil, nil
	}

	notif, err := g.repo.CreateNotification(ctx, repository.CreateNotificationParams{
		ID:              ulid.New(),
		RecipientUserID: params.RecipientUserID,
		Type:            params.Type,
		ActorUserID:     params.ActorUserID,
		ResourceName:    params.ResourceName,
	})
	if err != nil {
		return nil, err
	}

	view, err := g.BuildNotificationView(ctx, params.RecipientUserID, notif)
	if err != nil {
		return nil, err
	}

	g.publishCreated(ctx, view)
	return view, nil
}

// publishCreated is the fan-out companion to Create. Kept private because
// publish should only ever be paired with a successful DB write from Create.
func (g *gatewayImpl) publishCreated(ctx context.Context, view *NotificationView) {
	if view == nil || view.Notification == nil {
		return
	}
	n := view.Notification
	envelope := &eventv1.StreamEventsResponse{
		EventId:    n.ID.String(),
		CreateTime: timestamppb.New(n.CreateTime),
		Payload: &eventv1.StreamEventsResponse_NotificationCreated{
			NotificationCreated: &eventv1.NotificationCreatedEvent{
				Notification: ToProto(view),
			},
		},
	}
	if err := g.publisher.Publish(ctx, events.UserTopic(n.RecipientUserID), envelope); err != nil {
		slog.Warn("notification: publish failed",
			slog.String("recipient", n.RecipientUserID.String()),
			slog.String("err", err.Error()),
		)
	}
}
