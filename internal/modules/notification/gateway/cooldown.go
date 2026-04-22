package gateway

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"reverie.jp/reverie/internal/domain/entity"
	"reverie.jp/reverie/internal/platform/ulid"
)

// cooldownForType returns the per-(recipient, type, actor, resource) silent
// window. Zero disables the check. The intent is to prevent spam from rapid
// toggle actions (follow→unfollow→follow) that would otherwise delete+recreate
// the notification row and bypass client-side event_id dedup.
func cooldownForType(t entity.NotificationType) time.Duration {
	switch t {
	case entity.NotificationTypeUserFollowed:
		// Follow/unfollow toggles are the main abuse vector.
		return 1 * time.Hour
	case entity.NotificationTypeFollowingUserCallStarted:
		// resource_name is calls/{new_id} for every new call, so the cooldown
		// key is always unique — explicit zero documents the intent.
		return 0
	default:
		return 0
	}
}

// cooldownKey scopes the lock to the exact tuple that should dedupe. Different
// actors (or different resources) produce different keys, so independent
// events never block each other.
func cooldownKey(recipient ulid.ULID, t entity.NotificationType, actor *ulid.ULID, resource string) string {
	actorPart := ""
	if actor != nil {
		actorPart = actor.String()
	}
	return fmt.Sprintf("notif:cooldown:%s:%s:%s:%s", recipient.String(), t, actorPart, resource)
}

// acquireCooldown returns true when the caller should proceed with the
// notification. Returns false when the tuple was notified within the cooldown
// window (caller should silently skip). Redis errors fail open (proceed) so
// cooldown infra outages don't suppress legitimate notifications.
func (g *gatewayImpl) acquireCooldown(ctx context.Context, recipient ulid.ULID, t entity.NotificationType, actor *ulid.ULID, resource string) bool {
	ttl := cooldownForType(t)
	if ttl <= 0 {
		return true
	}
	ok, err := g.limiter.TryAcquire(ctx, cooldownKey(recipient, t, actor, resource), ttl)
	if err != nil {
		slog.Warn("notification cooldown: limiter error — failing open",
			slog.String("err", err.Error()),
		)
		return true
	}
	return ok
}
