package gateway

import (
	"context"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/modules/notification/repository"
	usergw "reverie.jp/reverie/internal/modules/user/gateway"
	"reverie.jp/reverie/internal/platform/events"
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

// FanOutParams describes a broadcast: every recipient gets a notification
// with the same actor / type / resource_name. Used by call_started → followers.
type FanOutParams struct {
	RecipientUserIDs []ulid.ULID
	Type             entity.NotificationType
	ActorUserID      *ulid.ULID
	ResourceName     string
}

type Gateway interface {
	// Create inserts a notification and publishes it on the recipient's user
	// topic. DB write is authoritative; publish failures are logged and
	// swallowed so event bus outages don't block the caller.
	Create(ctx context.Context, params CreateParams) (*NotificationView, error)

	// FanOutCreate batches the "same actor/type/resource to many recipients"
	// pattern: one INSERT for the DB writes, parallel publish to stream.
	// Publish runs in a detached goroutine so large fan-outs don't block the
	// caller. Follows stream-as-hint — DB write is authoritative.
	FanOutCreate(ctx context.Context, params FanOutParams) error

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
