package gateway

import (
	"context"
	"log/slog"

	"google.golang.org/protobuf/types/known/timestamppb"

	"reverie.jp/reverie/internal/domain/entity"
	eventv1 "reverie.jp/reverie/internal/gen/pb/event/v1"
	notificationv1 "reverie.jp/reverie/internal/gen/pb/notification/v1"
	"reverie.jp/reverie/internal/modules/notification/repository"
	useradapter "reverie.jp/reverie/internal/modules/user/adapter"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/events"
	"reverie.jp/reverie/internal/platform/resourcename"
	"reverie.jp/reverie/internal/platform/ulid"
)

// NotificationView is the composed read model returned to upstream modules.
// Notification itself lives in entity; the actor User is pulled via user
// gateway so we get the same view flags as other user surfaces.
type NotificationView struct {
	Notification *entity.Notification
	Actor        *usergw.UserView
}

type CreateParams struct {
	RecipientUserID ulid.ULID
	Type            entity.NotificationType
	ActorUserID     *ulid.ULID
	ResourceName    string
}

type Gateway interface {
	// Create inserts a notification and publishes it on the recipient's user
	// topic. DB write is authoritative; publish failures are logged and
	// swallowed so event bus outages don't block the caller.
	Create(ctx context.Context, params CreateParams) (*NotificationView, error)

	ListByRecipient(ctx context.Context, recipientID ulid.ULID, cursorID string, pageSize int32) ([]*NotificationView, error)
	MarkRead(ctx context.Context, recipientID ulid.ULID, ids []ulid.ULID) (int32, error)
	MarkAllRead(ctx context.Context, recipientID ulid.ULID) (int32, error)
	CountUnread(ctx context.Context, recipientID ulid.ULID) (int32, error)

	// DeleteByTypeActor removes notifications of a given type from a single
	// actor. Used by the follow usecase to clear "user_followed" notifications
	// on unfollow so a subsequent re-follow produces a fresh notification
	// instead of hitting the dedup index and returning the stale row.
	DeleteByTypeActor(ctx context.Context, recipientID ulid.ULID, notifType entity.NotificationType, actorID ulid.ULID) error

	BuildListNotificationViews(ctx context.Context, recipientID ulid.ULID, notifications []*entity.Notification) ([]*NotificationView, error)
	BuildNotificationView(ctx context.Context, recipientID ulid.ULID, n *entity.Notification) (*NotificationView, error)
}

type gatewayImpl struct {
	repo        repository.Repository
	userGateway usergw.Gateway
	publisher   events.Publisher
}

func New(repo repository.Repository, userGateway usergw.Gateway, publisher events.Publisher) Gateway {
	return &gatewayImpl{repo: repo, userGateway: userGateway, publisher: publisher}
}

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

func (g *gatewayImpl) ListByRecipient(ctx context.Context, recipientID ulid.ULID, cursorID string, pageSize int32) ([]*NotificationView, error) {
	notifs, err := g.repo.ListNotificationsByRecipient(ctx, recipientID, cursorID, pageSize)
	if err != nil {
		return nil, err
	}
	return g.BuildListNotificationViews(ctx, recipientID, notifs)
}

func (g *gatewayImpl) MarkRead(ctx context.Context, recipientID ulid.ULID, ids []ulid.ULID) (int32, error) {
	n, err := g.repo.MarkNotificationsRead(ctx, recipientID, ids)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}

func (g *gatewayImpl) MarkAllRead(ctx context.Context, recipientID ulid.ULID) (int32, error) {
	n, err := g.repo.MarkAllNotificationsRead(ctx, recipientID)
	if err != nil {
		return 0, err
	}
	return int32(n), nil
}

func (g *gatewayImpl) CountUnread(ctx context.Context, recipientID ulid.ULID) (int32, error) {
	return g.repo.CountUnreadNotifications(ctx, recipientID)
}

func (g *gatewayImpl) DeleteByTypeActor(ctx context.Context, recipientID ulid.ULID, notifType entity.NotificationType, actorID ulid.ULID) error {
	_, err := g.repo.DeleteNotificationsByTypeActor(ctx, recipientID, notifType, actorID)
	return err
}

func (g *gatewayImpl) BuildNotificationView(ctx context.Context, recipientID ulid.ULID, n *entity.Notification) (*NotificationView, error) {
	views, err := g.BuildListNotificationViews(ctx, recipientID, []*entity.Notification{n})
	if err != nil {
		return nil, err
	}
	if len(views) == 0 {
		return nil, nil
	}
	return views[0], nil
}

func (g *gatewayImpl) BuildListNotificationViews(ctx context.Context, recipientID ulid.ULID, notifications []*entity.Notification) ([]*NotificationView, error) {
	if len(notifications) == 0 {
		return []*NotificationView{}, nil
	}
	actorIDs := make([]ulid.ULID, 0, len(notifications))
	seen := make(map[ulid.ULID]struct{}, len(notifications))
	for _, n := range notifications {
		if n == nil || n.ActorUserID == nil {
			continue
		}
		if _, ok := seen[*n.ActorUserID]; ok {
			continue
		}
		seen[*n.ActorUserID] = struct{}{}
		actorIDs = append(actorIDs, *n.ActorUserID)
	}

	actorsByID := map[ulid.ULID]*usergw.UserView{}
	if len(actorIDs) > 0 {
		actors, err := g.userGateway.BuildListUserViews(ctx, recipientID, actorIDs)
		if err != nil {
			return nil, err
		}
		for _, v := range actors {
			if v != nil && v.User != nil {
				actorsByID[v.User.ID] = v
			}
		}
	}

	out := make([]*NotificationView, 0, len(notifications))
	for _, n := range notifications {
		if n == nil {
			continue
		}
		var actor *usergw.UserView
		if n.ActorUserID != nil {
			actor = actorsByID[*n.ActorUserID]
		}
		out = append(out, &NotificationView{Notification: n, Actor: actor})
	}
	return out, nil
}

// ToProto converts a view into the proto Notification shared between the
// ListNotifications response and the event envelope payload. Kept here so
// both usages produce identical output.
func ToProto(view *NotificationView) *notificationv1.Notification {
	if view == nil || view.Notification == nil {
		return nil
	}
	n := view.Notification
	pb := &notificationv1.Notification{
		Type:         toProtoType(n.Type),
		ResourceName: n.ResourceName,
		CreateTime:   timestamppb.New(n.CreateTime),
	}
	if n.ReadTime != nil {
		pb.ReadTime = timestamppb.New(*n.ReadTime)
	}
	if view.Actor != nil && view.Actor.User != nil {
		pb.Name = resourcename.FormatNotification(view.Actor.User.CustomID, n.ID)
		pb.Actor = useradapter.ToUser(view.Actor)
	} else {
		pb.Name = "notifications/" + n.ID.String()
	}
	return pb
}

func toProtoType(t entity.NotificationType) notificationv1.NotificationType {
	switch t {
	case entity.NotificationTypeUserFollowed:
		return notificationv1.NotificationType_NOTIFICATION_TYPE_USER_FOLLOWED
	case entity.NotificationTypeFollowingUserCallStarted:
		return notificationv1.NotificationType_NOTIFICATION_TYPE_FOLLOWING_USER_CALL_STARTED
	default:
		return notificationv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED
	}
}

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
