package gateway

import (
	"context"
	"log/slog"
	"sync"

	"google.golang.org/protobuf/types/known/timestamppb"

	"reverie.jp/reverie/internal/domain/entity"
	eventv1 "reverie.jp/reverie/internal/gen/pb/event/v1"
	"reverie.jp/reverie/internal/modules/notification/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/events"
	"reverie.jp/reverie/internal/platform/ulid"
)

// fanOutPublishConcurrency bounds how many Redis publishes run in parallel
// for a single fan-out. 20 balances throughput with connection pressure on
// the events bus.
const fanOutPublishConcurrency = 20

func (g *gatewayImpl) FanOutCreate(ctx context.Context, params FanOutParams) error {
	if len(params.RecipientUserIDs) == 0 {
		return nil
	}

	// Drop self-notification (actor happens to be in the recipient list).
	recipients := make([]ulid.ULID, 0, len(params.RecipientUserIDs))
	seen := make(map[ulid.ULID]struct{}, len(params.RecipientUserIDs))
	for _, rid := range params.RecipientUserIDs {
		if rid.IsZero() {
			continue
		}
		if params.ActorUserID != nil && *params.ActorUserID == rid {
			continue
		}
		if _, dup := seen[rid]; dup {
			continue
		}
		seen[rid] = struct{}{}
		recipients = append(recipients, rid)
	}
	if len(recipients) == 0 {
		return nil
	}

	ids := make([]ulid.ULID, len(recipients))
	for i := range recipients {
		ids[i] = ulid.New()
	}

	inserted, err := g.repo.CreateFanOutNotifications(ctx, repository.CreateFanOutNotificationsParams{
		IDs:          ids,
		RecipientIDs: recipients,
		Type:         params.Type,
		ActorUserID:  params.ActorUserID,
		ResourceName: params.ResourceName,
	})
	if err != nil {
		return err
	}
	if len(inserted) == 0 {
		return nil
	}

	// Build the actor view once. Fan-out publishes deliberately omit
	// per-recipient relationship flags (no IsFollowing etc.) — the view is
	// shown in a toast, and the bell's ListNotifications path re-builds with
	// full flags for anyone who opens it.
	var actor *usergw.UserView
	if params.ActorUserID != nil {
		actor, err = g.userGateway.BuildUserView(ctx, ulid.ULID{}, *params.ActorUserID)
		if err != nil {
			// DB rows are already in — log and skip publish. Recipients will
			// see the notification on next ListNotifications / page refresh.
			slog.Warn("notification fan-out: build actor view failed",
				slog.String("actor", params.ActorUserID.String()),
				slog.String("err", err.Error()),
			)
			return nil
		}
	}

	// Detach from the caller's context so a quick RPC return doesn't cancel
	// pending publishes. Values (trace IDs etc.) are preserved.
	go g.publishFanOut(context.WithoutCancel(ctx), inserted, actor)
	return nil
}

func (g *gatewayImpl) publishFanOut(ctx context.Context, notifications []*entity.Notification, actor *usergw.UserView) {
	sem := make(chan struct{}, fanOutPublishConcurrency)
	var wg sync.WaitGroup
	for _, n := range notifications {
		if n == nil {
			continue
		}
		wg.Add(1)
		sem <- struct{}{}
		go func(n *entity.Notification) {
			defer wg.Done()
			defer func() { <-sem }()
			view := &NotificationView{Notification: n, Actor: actor}
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
				slog.Warn("notification fan-out: publish failed",
					slog.String("recipient", n.RecipientUserID.String()),
					slog.String("err", err.Error()),
				)
			}
		}(n)
	}
	wg.Wait()
}
